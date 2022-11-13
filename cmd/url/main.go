package main

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	app "github.com/kormiltsev/url-shorter/app"
	storage "github.com/kormiltsev/url-shorter/db"
)

const keyServerAddr = "losalhost"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read body: %s\n", err)
	}

	log.Printf("%s: got / request. Data: %s\n",
		ctx.Value(keyServerAddr),
		body)
	var newrow storage.Baserow
	if r.Method == "POST" {
		newrow = storage.Baserow{
			Surl: app.GetRandomString(),
			Lurl: string(body),
		}
		err := storage.QueryPOSTorSelect(&newrow)
		if err != nil {
			log.Printf("%s, POST Query error. %s", err, newrow.Lurl)
		}
		io.WriteString(w, newrow.Surl)
	} else if r.Method == "GET" {
		newrow = storage.Baserow{
			Surl: string(body),
		}
		err := storage.QueryGetURL(&newrow)
		if err != nil {
			log.Printf("%s in query GET Surl %s", err, string(body))
			//anwser 404
		} else {
			io.WriteString(w, newrow.Lurl)
		}
	} else {
		http.Error(w, "Invalid request method.", 405)
	}

}

func main() {
	// connect to postgres
	storage.StartConnection()
	defer storage.DbClose()
	qty, err := storage.RowsQuantity()
	if err != nil {
		log.Printf("%s. Postgres cant query rows quantity\n", err)
	} else {
		log.Println("Total rows in postgres = ", qty)
	}
	// ================
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot) //.Methods("GET")

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server closed\n")
	} else if err != nil {
		log.Fatal("error listening for server: %s\n", err)
	}

}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

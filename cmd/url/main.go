package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	app "github.com/kormiltsev/url-shorter/app"
	storage "github.com/kormiltsev/url-shorter/db"
	//db "github.com/kormiltsev/url-shorter/db"
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
			Lurl: string(body),
		}
		err := storage.QueryGetURL(&newrow)
		if err != nil {
			newrow.Surl = GetNewSurl(string(body))
			storage.InsertNewRow(&newrow)
		}
		io.WriteString(w, newrow.Surl)
	} else {
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
	}

}
func getInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: got /info request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "GET_POST\n")
}

func main() {
	// connect to postgres
	err := storage.StartConnection()
	if err != nil {
		log.Printf("%s. Postgres not connected\n", err)
	}
	err = storage.CreateTable()
	if err != nil {
		log.Printf("%s. Postgres cant create table\n", err)
	}
	qty, err := storage.RowsQuantity()
	if err != nil {
		log.Printf("%s. Postgres cant query rows quantity\n", err)
	} else {
		log.Println("Total rows in postgres = ", qty)
	}
	// ================
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/help", getInfo)

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

func GetNewSurl(u string) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	var url_len = 10
	return app.GetRandomString(url_len, letters)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	app "github.com/kormiltsev/url-shorter/app"
	storage "github.com/kormiltsev/url-shorter/db"
)

const keyServerAddr = "losalhost"

var port string
var postgres bool
var filedir string
var conf = [4]string{}

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
		err := storage.Router(postgres, &newrow, r.Method)
		if err != nil {
			log.Printf("%s, POST Query error. %s", err, newrow.Lurl)
		}
		io.WriteString(w, newrow.Surl)
	} else if r.Method == "GET" {
		newrow = storage.Baserow{
			Surl: string(body),
		}
		err := storage.Router(postgres, &newrow, r.Method)
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

func Flags() {
	// in case of .env
	godotenv.Load()
	conf[0] = os.Getenv("ADR")
	conf[1] = os.Getenv("USR")
	conf[2] = os.Getenv("PASS")
	conf[3] = os.Getenv("DB")

	por := os.Getenv("PORT")
	dir := os.Getenv("LOCALFILE")
	// in case of flags thru terminal
	{
		a := flag.String("adr", conf[0], "a string")
		u := flag.String("usr", conf[1], "a string")
		p := flag.String("pas", conf[2], "a string")
		d := flag.String("db", conf[3], "a string")
		por := flag.String("port", por, "request url via this port")
		i := flag.Bool("pg", false, "db = postgres. default (false): internal file")
		f := flag.String("files", dir, "for internal db files")
		flag.Parse()
		conf[0] = *a
		conf[1] = *u
		conf[2] = *p
		conf[3] = *d

		port = *por
		postgres = *i
		filedir = *f
	}
	// in case of any error use default
	if conf[0] == "" || conf[1] == "" || conf[2] == "" || conf[3] == "" {
		conf[0] = "localhost:5432"
		conf[1] = "postgres"
		conf[2] = "root"
		conf[3] = "postgres"
	}
}
func main() {
	// get env and flags
	Flags()
	// connect to postgres
	if postgres {
		storage.StartConnection(conf)
		defer storage.DbClose()
		qty, err := storage.RowsQuantity()
		if err != nil {
			log.Printf("%s. Postgres cant query rows quantity\n", err)
		} else {
			log.Println("Total rows in postgres = ", qty)
		}
	} else {
		err := storage.ConnectFile(filedir)
		if err != nil {
			return
		}
	}
	// ================
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot) //.Methods("GET")

	ctx := context.Background()
	server := &http.Server{
		Addr:    port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	err := server.ListenAndServe()
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

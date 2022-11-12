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
)

const keyServerAddr = "losalhost"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request. Data: %s\n",
		ctx.Value(keyServerAddr),
		body)
	if r.Method == "POST" {
		io.WriteString(w, randAnswer(string(body))+"\n")
	} else {
		io.WriteString(w, "answer!\n")
	}

}
func getInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: got /info request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "GET_POST\n")
}
func main() {
	// connect to postgres

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
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server closed\n")
	} else if err != nil {
		log.Fatal("error listening for server: %s\n", err)
	}

}

func randAnswer(u string) string {
	// const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	// var url_len = 10
	//return GetRandomString(url_len, letters)
	return "result " + u
}

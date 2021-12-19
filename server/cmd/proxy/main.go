package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ihaxolotl/webproxy/internal/api"
	"github.com/ihaxolotl/webproxy/internal/data"
)

const APIAddr = ":8888"

func main() {
	db := data.New()
	if err := db.Setup(); err != nil {
		log.Fatal(err)
	}

	m := mux.NewRouter()
	m.StrictSlash(true)

	ctx := api.Context{Database: db}

	for _, rt := range api.APIRoutes {
		m.Path(rt.URL).
			Name(rt.Name).
			Methods(rt.Method).
			Handler(rt.Handler(ctx))
	}

	s := &http.Server{
		Addr:         APIAddr,
		Handler:      m,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Fatal(s.ListenAndServe())
}

package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ihaxolotl/webproxy/internal/api"
	"github.com/ihaxolotl/webproxy/internal/data"
	"github.com/ihaxolotl/webproxy/internal/proxy"
)

const ProxyListenerAddr = ":8080"
const APIAddr = ":8888"

func main() {
	db, err := data.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		listener, err := net.Listen("tcp", ProxyListenerAddr)
		if err != nil {
			log.Fatal(err)
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}

			go proxy.HandleRequest(conn, db)
		}
	}()

	r := mux.NewRouter()

	for _, rt := range api.APIRoutes {
		r.HandleFunc(rt.URL, rt.Handler).Name(rt.Name).Methods(rt.Method)
	}

	s := &http.Server{
		Addr:         APIAddr,
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Fatal(s.ListenAndServe())
}

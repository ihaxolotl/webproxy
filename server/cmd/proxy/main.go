package main

import (
	"log"
	"net"

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
}

package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

const ProxyListenerAddr = ":8080"
const APIAddr = ":8888"
const ReadBufferSize = 0xffff

func handleProxyRequest(conn net.Conn) {
	var (
		clientRequest  *Buffer
		proxyRequest   *Buffer
		serverResponse *Buffer
		hostname       string
		err            error
	)

	// Read a request from the client.
	clientRequest = NewBuffer()
	if err = clientRequest.Recv(conn); err != nil {
		if err != io.EOF {
			log.Fatalln("read error: ", err)
		}

		log.Println(err)
	}
	defer conn.Close()

	// HACK: Parse the the request to get the hostname.
	dummy := readRequest(clientRequest.Buffer(), clientRequest.Size())

	// FIXME: Discard all CONNECT requests
	// Let's not deal with HTTPS yet.
	if dummy.Method == http.MethodConnect {
		return
	}

	// Ensure that the hostname format is always host:port.
	hostname = dummy.Host
	if dummy.URL.Port() == "" {
		hostname = hostname + ":80"
	}

	// Connect to the target server.
	proxyConn, err := net.Dial("tcp", hostname)
	if err != nil {
		log.Fatal(err)
	}
	defer proxyConn.Close()

	// HACK: Replace the proxy headers in the request.
	filters := []filter{
		{Find: "Proxy-Connection:", Replace: "Connection:"},
		{Find: dummy.URL.Scheme + "://" + dummy.URL.Host, Replace: ""},
	}

	if proxyRequest, err = parseProxyRequest(clientRequest, filters); err != nil {
		log.Fatalln("parse error: ", err)
	}

	// Proxy the request to its destination.
	if err = proxyRequest.Send(proxyConn); err != nil {
		log.Fatalln("write error: ", err)
	}

	// Read the server response and send it back to the client connection.
	serverResponse = NewBuffer()
	if err = serverResponse.Recvall(proxyConn); err != nil {
		log.Fatalln("read error: ", err)
	}

	if err = serverResponse.Send(conn); err != nil {
		log.Fatalln("write error: ", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ProxyListenerAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleProxyRequest(conn)
	}
}

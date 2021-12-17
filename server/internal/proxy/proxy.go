package proxy

import (
	"database/sql"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ihaxolotl/webproxy/internal/buffer"
	"github.com/ihaxolotl/webproxy/internal/data"
)

func HandleRequest(conn net.Conn, db *sql.DB) {
	var (
		clientRequest  *buffer.Buffer
		proxyRequest   *buffer.Buffer
		serverResponse *buffer.Buffer
		hostname       string
		err            error
	)

	// Read a request from the client.
	clientRequest = buffer.NewBuffer()
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

	//
	// TODO(Brett): Stall the request here
	//

	requestEntry := data.Request{
		ID:        uuid.New().String(),
		ProjectID: "NoneYet",
		Method:    dummy.Method,
		Domain:    dummy.URL.Host,
		IPAddr:    dummy.URL.Host,
		Length:    int64(len(proxyRequest.Buffer())),
		Edited:    false,
		Comment:   "",
		Raw:       string(proxyRequest.Buffer()),
	}

	if _, err = data.InsertRequest(db, &requestEntry); err != nil {
		log.Fatalln("database: ", err)
	}

	// Proxy the request to its destination.
	if err = proxyRequest.Send(proxyConn); err != nil {
		log.Fatalln("write error: ", err)
	}

	// Read the server response and send it back to the client connection.
	serverResponse = buffer.NewBuffer()
	if err = serverResponse.Recvall(proxyConn); err != nil {
		log.Fatalln("read error: ", err)
	}

	responseEntry := data.Response{
		ID:        uuid.New().String(),
		ProjectID: "NoneYet",
		Status:    0,
		Length:    int64(len(serverResponse.Buffer())),
		Edited:    false,
		Timestamp: time.Now(),
		Mimetype:  "NoneYet",
		Comment:   "",
		Raw:       string(serverResponse.Buffer()),
	}

	if _, err = data.InsertResponse(db, &responseEntry); err != nil {
		log.Fatalln("database: ", err)
	}

	if err = serverResponse.Send(conn); err != nil {
		log.Fatalln("write error: ", err)
	}
}

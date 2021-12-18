package proxy

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ihaxolotl/webproxy/internal/buffer"
	"github.com/ihaxolotl/webproxy/internal/data"
)

const DefaultPort = 8080

var ErrUnknownProxyCmd = errors.New("unknown proxy command")

type ProxyCmdType byte

const (
	ProxyCmdUnknown ProxyCmdType = iota
	ProxyCmdStart
	ProxyCmdStop
	ProxyCmdStall
	ProxyCmdForward
	ProxyCmdDrop
)

var proxyCmdTypes = map[ProxyCmdType]string{
	ProxyCmdUnknown: "ProxyCmdUnknown",
	ProxyCmdStart:   "ProxyCmdStart",
	ProxyCmdStop:    "ProxyCmdStop",
	ProxyCmdStall:   "ProxyCmdStall",
	ProxyCmdForward: "ProxyCmdForward",
	ProxyCmdDrop:    "ProxyCmdDrop",
}

func (m ProxyCmdType) String() string {
	return proxyCmdTypes[m]
}

type ProxyCmd struct {
	Type ProxyCmdType `json:"type"`
	Data []byte       `json:"data"`
}

// Options represents the configuration object for the proxy.
type Options struct {
	ListenPort      int  // Proxy listener port
	InterceptClient bool // Intercept client HTTP requests
	InterceptServer bool // Intercept server HTTP responses
	Stall           bool // Stall enables stalling requests/responses
}

type Proxy struct {
	db   *sql.DB
	conn *websocket.Conn
	cmd  chan ProxyCmd
}

func New(db *sql.DB, conn *websocket.Conn, cmd chan ProxyCmd) *Proxy {
	return &Proxy{db, conn, cmd}
}

func (proxy *Proxy) Spawn() {
	var (
		listener net.Listener
		opts     Options
		err      error
	)

	opts = Options{
		ListenPort:      DefaultPort,
		InterceptClient: true,
		InterceptServer: true,
		Stall:           true,
	}

	listener, err = net.Listen("tcp", fmt.Sprintf(":%d", opts.ListenPort))
	if err != nil {
		log.Fatal(err)
	}

	for {
		var conn net.Conn

		conn, err = listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		proxy.HandleRequest(conn, &opts)
	}
}

// HandleRequest handles requests and response by acting as a middle-man.
// Requests are received from the client and forwarded to their destination.
func (proxy *Proxy) HandleRequest(conn net.Conn, opts *Options) {
	var (
		clientRequest  *buffer.Buffer
		proxyRequest   *buffer.Buffer
		serverResponse *buffer.Buffer
		hostname       string
		err            error
	)

	defer conn.Close()

	// Read a request from the client.
	clientRequest = buffer.NewBuffer()
	if err = clientRequest.Recv(conn); err != nil {
		if err != io.EOF {
			log.Fatalln("read error: ", err)
		}

		log.Println(err)
	}

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

	// Stall requests
	if opts.InterceptClient && opts.Stall {
		if err = proxy.conn.WriteMessage(
			websocket.TextMessage, proxyRequest.Buffer(),
		); err != nil {
			log.Fatal(err)
		}

		cmd := <-proxy.cmd
		if cmd.Type == ProxyCmdDrop {
			return
		}
	}

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

	if _, err = data.InsertRequest(proxy.db, &requestEntry); err != nil {
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

	// Stall responses
	if opts.InterceptServer && opts.Stall {
		if err = proxy.conn.WriteMessage(
			websocket.TextMessage, serverResponse.Buffer(),
		); err != nil {
			log.Fatal(err)
		}

		cmd := <-proxy.cmd
		if cmd.Type == ProxyCmdDrop {
			return
		}
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

	if _, err = data.InsertResponse(proxy.db, &responseEntry); err != nil {
		log.Fatalln("database: ", err)
	}

	if err = serverResponse.Send(conn); err != nil {
		log.Fatalln("write error: ", err)
	}
}

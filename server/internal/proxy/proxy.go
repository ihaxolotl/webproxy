package proxy

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

// Options represents the configuration object for the proxy.
type Options struct {
	ListenPort      int  // Proxy listener port
	InterceptClient bool // Intercept client HTTP requests
	InterceptServer bool // Intercept server HTTP responses
	Stall           bool // Stall enables stalling requests/responses
}

type Proxy struct {
	db     *sql.DB
	conn   *websocket.Conn
	cmd    chan ProxyCmd
	intcmd chan ProxyCmd
}

// New allocates memory for and returns a Proxy.
func New(db *sql.DB, conn *websocket.Conn, cmd chan ProxyCmd) *Proxy {
	return &Proxy{
		db:     db,
		conn:   conn,
		cmd:    cmd,
		intcmd: make(chan ProxyCmd),
	}
}

// Spawn starts a new TCP proxy listener and accepts requests from the client.
func (proxy *Proxy) Spawn() {
	var (
		listener net.Listener
		opts     Options
		err      error
	)

	// Set default options for interception
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

	// FIXME: Would a mutex be required for accessing this critical section?
	go func() {
		for {
			cmd := <-proxy.cmd
			switch cmd.Type {
			case ProxyCmdStart:
				opts.Stall = true
				fmt.Printf("Stall: on\n")
			case ProxyCmdStop:
				opts.Stall = false
				fmt.Printf("Stall: off\n")
			case ProxyCmdForward, ProxyCmdDrop:
				proxy.intcmd <- cmd
			default:
			}
		}
	}()

	for {
		var conn net.Conn

		conn, err = listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		proxy.HandleRequest(conn, &opts)
	}
}

// stall takes intercepted data and sends it to the client's WebSocket connection and blocks
// until a command is received. If the command type if ProxyCmdForward, the original request
// is compared to the command data. If the two buffers match, the edited flag will be set.
// If the command type is ProxyCmdDrop, return an error.
func (proxy *Proxy) stall(stalled *buffer.Buffer, edited *bool) (*buffer.Buffer, error) {
	var (
		cmd       ProxyCmd
		msg       ProxyCmd
		payload   []byte
		forwarded *buffer.Buffer
		err       error
	)

	msg = ProxyCmd{
		Type: ProxyCmdStall,
		Data: string(stalled.Buffer()),
	}

	if payload, err = json.Marshal(&msg); err != nil {
		return nil, err
	}

	if err = proxy.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return nil, err
	}

	cmd = <-proxy.intcmd
	if cmd.Type == ProxyCmdForward {
		if forwarded = buffer.NewBufferFrom([]byte(cmd.Data), len(cmd.Data)); forwarded == nil {
			return nil, ErrNilBuffer
		}

		if bytes.Compare(forwarded.Buffer(), stalled.Buffer()) != 0 {
			*edited = true
		}

		return forwarded, nil
	}

	return nil, ErrDropped
}

// HandleRequest handles requests and response by acting as a middle-man.
// Requests are received from the client and forwarded to their destination.
func (proxy *Proxy) HandleRequest(conn net.Conn, opts *Options) {
	var (
		clientRequest  *buffer.Buffer
		proxyRequest   *buffer.Buffer
		serverResponse *buffer.Buffer
		hostname       string
		requestEdited  bool
		responseEdited bool
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
		if proxyRequest, err = proxy.stall(proxyRequest, &requestEdited); err != nil {
			if err != ErrDropped {
				log.Fatal(err)
			}

			log.Println(err)
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
		Edited:    requestEdited,
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
		if serverResponse, err = proxy.stall(serverResponse, &responseEdited); err != nil {
			if err != ErrDropped {
				return
			}
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

package proxy

import (
	"bytes"
	"encoding/json"
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
	"github.com/ihaxolotl/webproxy/internal/data/requests"
	"github.com/ihaxolotl/webproxy/internal/data/responses"
)

var ErrTmpNoOptions = errors.New("OPTIONS requests not allowed (yet)")

// DefaultPort for testing. The user should soon be able to set whichever
// port they want to listen on in their project settings.
const DefaultPort = 8080

// Options represents the configuration object for the proxy.
type Options struct {
	ListenPort      int  // Proxy listener port
	InterceptClient bool // Intercept client HTTP requests
	InterceptServer bool // Intercept server HTTP responses
	Stall           bool // Stall enables stalling requests/responses
}

// Proxy is an intercepting proxy server.
type Proxy struct {
	projectId string          // Unique ID of the project. NOTE: This may be moved.
	db        *data.Database  // Database connection
	conn      *websocket.Conn // Client WebSocket connection
	cmd       chan ProxyCmd   // Command queue channel
	intcmd    chan ProxyCmd   // Intercept behaviour command queue channel
}

// New allocates memory for and returns a Proxy.
func New(
	projectId string,
	db *data.Database,
	conn *websocket.Conn,
	cmd chan ProxyCmd,
) *Proxy {
	return &Proxy{
		projectId: projectId,
		db:        db,
		conn:      conn,
		cmd:       cmd,
		intcmd:    make(chan ProxyCmd),
	}
}

// Spawn creates a new TCP proxy listener and accepts connections from the client.
// The connections accepted by the listener will have requests and responses that
// can be stalled and modified at the control panel. The listener handles all
// connections synchronously to avoid race conditions.
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
		Stall:           false,
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
		defer conn.Close()

		if err = proxy.HandleRequest(conn, &opts); err != nil {
			log.Println(err)
		}
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

type httpdata struct {
	Request          *http.Request
	Response         *http.Response
	RawRequest       *buffer.Buffer
	RawResponse      *buffer.Buffer
	Elapsed          time.Duration
	RequestTime      time.Time
	ResponseTime     time.Time
	IsRequestEdited  bool
	IsResponseEdited bool
}

// commit inserts the data contained in the passed httpdata struct into the
// appropriate tables in the database.
func (proxy *Proxy) commit(d *httpdata) error {
	var (
		requestId      string
		responseId     string
		requestRecord  requests.Request
		responseRecord responses.Response
		err            error
	)

	requestId = uuid.New().String()
	responseId = uuid.New().String()

	requestRecord = requests.Request{
		ID:         requestId,
		ProjectID:  proxy.projectId,
		ResponseID: responseId,
		Method:     d.Request.Method,
		Domain:     d.Request.URL.Host,
		IPAddr:     d.Request.URL.Host, // TODO(Brett) Record the IP address of hosts
		URL:        d.Request.URL.RequestURI(),
		Length:     int64(d.RawRequest.Size()),
		Edited:     d.IsRequestEdited,
		Timestamp:  d.RequestTime,
		Comment:    "", // TODO(Brett): Implement comments
		Raw:        string(d.RawRequest.Buffer()),
	}

	if _, err = proxy.db.Requests.Insert(&requestRecord); err != nil {
		return err
	}

	responseRecord = responses.Response{
		ID:        responseId,
		ProjectID: proxy.projectId,
		RequestID: requestId,
		Status:    int16(d.Response.StatusCode),
		Length:    int64(d.RawResponse.Size()),
		Edited:    d.IsResponseEdited,
		Elapsed:   int64(d.Elapsed),
		Timestamp: d.ResponseTime,
		Mimetype:  "", // TODO(Brett): Record response body mime-types
		Comment:   "", // TODO(Brett): Implement comments
		Raw:       string(d.RawResponse.Buffer()),
	}

	if _, err = proxy.db.Responses.Insert(&responseRecord); err != nil {
		return err
	}

	return err
}

// HandleRequest handles requests and response by acting as a middle-man.
// Requests are received from the client and forwarded to their destination.
func (proxy *Proxy) HandleRequest(conn net.Conn, opts *Options) error {
	var (
		clientRequest  *buffer.Buffer
		proxyRequest   *buffer.Buffer
		serverResponse *buffer.Buffer
		httpRequest    *http.Request
		dbdata         httpdata
		hostname       string
		timer          time.Time
		err            error
	)

	dbdata = httpdata{}

	// Read client request
	clientRequest = buffer.NewBuffer()
	if err = clientRequest.Recv(conn); err != nil {
		if err != io.EOF {
			return err
		}

		return nil
	}

	httpRequest = readRequest(clientRequest)

	// FIXME: Discard all CONNECT requests
	// Let's not deal with HTTPS yet.
	if httpRequest.Method == http.MethodConnect {
		return ErrTmpNoOptions
	}

	// Send the client's request to the target server.
	if proxyRequest, err = parseProxyRequest(clientRequest, httpRequest); err != nil {
		return err
	}

	// Stall requests
	if opts.InterceptClient && opts.Stall {
		proxyRequest, err = proxy.stall(proxyRequest, &dbdata.IsRequestEdited)
		if err != nil {
			if err != ErrDropped {
				return err
			}

			return nil
		}
	}

	dbdata.Request = httpRequest
	dbdata.RawRequest = proxyRequest

	// Ensure that the hostname format is always host:port.
	hostname = httpRequest.Host
	if httpRequest.URL.Port() == "" {
		hostname = hostname + ":80"
	}

	// Connect to the target server.
	proxyConn, err := net.Dial("tcp", hostname)
	if err != nil {
		return err
	}
	defer proxyConn.Close()

	// Proxy the request to its destination.
	if err = proxyRequest.Send(proxyConn); err != nil {
		return err
	}

	dbdata.RequestTime = time.Now()
	timer = time.Now()

	// Read the server response and send it back to the client connection.
	serverResponse = buffer.NewBuffer()
	if err = serverResponse.Recvall(proxyConn); err != nil {
		return err
	}

	dbdata.RawResponse = serverResponse
	dbdata.ResponseTime = time.Now()
	dbdata.Elapsed = dbdata.ResponseTime.Sub(timer)
	dbdata.Response = readResponse(httpRequest, serverResponse)

	// Stall responses
	if opts.InterceptServer && opts.Stall {
		serverResponse, err = proxy.stall(serverResponse, &dbdata.IsResponseEdited)
		if err != nil {
			if err != ErrDropped {
				return err
			}

			return nil
		}
	}

	if err = proxy.commit(&dbdata); err != nil {
		return nil
	}

	return serverResponse.Send(conn)
}

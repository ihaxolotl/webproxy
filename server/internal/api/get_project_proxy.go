package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ihaxolotl/webproxy/internal/proxy"
)

func HandleProxy(ctx Context, conn *websocket.Conn, cmd chan proxy.ProxyCmd) error {
	for {
		var (
			raw []byte
			msg proxy.ProxyCmd
			err error
		)

		_, raw, err = conn.ReadMessage()
		if err != nil {
			return err
		}

		if err = json.Unmarshal(raw, &msg); err != nil {
			return err
		}

		if err = msg.Validate(); err != nil {
			return err
		}

		cmd <- msg
	}
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{}

// GetProjectProxyRoute is an endpoint for connecting to the intercept proxy
// for the project. The endpoint will first check if the projectId passed as a
// URL variable corresponds to an existing project in the database. If the project
// exists, the endpoint will upgrade the connection to a WebSocket and will now
// receive messages from the client to control the proxy.
func GetProjectProxyRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			conn      *websocket.Conn
			projectId string
			vars      map[string]string
			cmd       chan proxy.ProxyCmd
			prox      *proxy.Proxy
			err       error
		)

		vars = mux.Vars(r)
		projectId = vars["projectId"]

		if _, err = ctx.Database.Projects.FetchById(projectId); err != nil {
			ctx.JSON(&rw, http.StatusNotFound, JSON{"err": err.Error()})
			return
		}

		if conn, err = upgrader.Upgrade(rw, r, nil); err != nil {
			ctx.JSON(&rw, http.StatusInternalServerError, JSON{"err": err.Error()})
			return
		}
		defer conn.Close()

		cmd = make(chan proxy.ProxyCmd)

		// Spawn a new intercept proxy
		prox = proxy.New(projectId, ctx.Database, conn, cmd)
		go prox.Spawn()

		if err = HandleProxy(ctx, conn, cmd); err != nil {
			log.Println(err)
			conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		}
	}
}

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ihaxolotl/webproxy/internal/data/requests"
)

// GetRequestByIdRoute is an endpoint that fetches a request matching a requestId
// passed as a URL variable. If the request does not exist, a status 404 is sent.
// If the request is found, it will be returned in full.
func GetRequestByIdRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			vars      map[string]string
			requestId string
			req       *requests.Request
			err       error
		)

		vars = mux.Vars(r)
		requestId = vars["requestId"]

		if req, err = ctx.Database.Requests.FetchById(requestId); err != nil {
			ctx.JSON(&rw, http.StatusNotFound, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusOK, JSON{"request": req})
	}
}

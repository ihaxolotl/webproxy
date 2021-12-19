package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ihaxolotl/webproxy/internal/data/responses"
)

// GetResponseByIdRoute is an endpoint that fetches a response matching a responseId
// passed as a URL variable. If the response does not exist, a status 404 is sent.
// If the response is found, it will be returned in full.
func GetResponseByIdRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			vars       map[string]string
			responseId string
			res        *responses.Response
			err        error
		)

		vars = mux.Vars(r)
		responseId = vars["responseId"]

		if res, err = ctx.Database.Responses.FetchById(responseId); err != nil {
			ctx.JSON(&rw, http.StatusNotFound, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusOK, JSON{"response": res})
	}
}

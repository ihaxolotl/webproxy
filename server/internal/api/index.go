package api

import (
	"net/http"
)

func IndexRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx.JSON(&rw, http.StatusOK, JSON{"msg": "Hello world!"})
	}
}

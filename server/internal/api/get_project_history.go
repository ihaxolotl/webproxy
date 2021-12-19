package api

import "net/http"

func GetProjectHistoryRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx.JSON(&rw, http.StatusNotImplemented, JSON{"msg": "Not implemented"})
	}
}

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ihaxolotl/webproxy/internal/data/history"
)

func GetProjectHistoryRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			hist      []history.HistoryEntry
			vars      map[string]string
			projectId string
			err       error
		)

		vars = mux.Vars(r)
		projectId = vars["projectId"]

		if hist, err = ctx.Database.History.Fetch(projectId); err != nil {
			ctx.JSON(&rw, http.StatusNotFound, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusOK, JSON{"history": hist})
	}
}

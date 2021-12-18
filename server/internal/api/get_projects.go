package api

import (
	"net/http"

	"github.com/ihaxolotl/webproxy/internal/data"
)

// GetProjectsRoute is an endpoint fetching all projects.
func GetProjectsRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			projects []data.Project
			err      error
		)

		projects, err = data.GetProjects(ctx.Database)
		if err != nil {
			ctx.JSON(&rw, http.StatusInternalServerError, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusOK, JSON{"projects": projects})
	}
}

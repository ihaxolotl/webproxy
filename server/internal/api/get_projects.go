package api

import (
	"net/http"

	"github.com/ihaxolotl/webproxy/internal/data/projects"
)

// GetProjectsRoute is an endpoint fetching all projects.
func GetProjectsRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			projs []projects.Project
			err   error
		)

		projs, err = ctx.Database.Projects.Fetch()
		if err != nil {
			ctx.JSON(&rw, http.StatusInternalServerError, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusOK, JSON{"projects": projs})
	}
}

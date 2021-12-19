package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ihaxolotl/webproxy/internal/data/projects"
)

// GetProjectsByIdRoute is an endpoint for fetching a project by its id.
func GetProjectsByIdRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			projectId string
			proj      *projects.Project
			vars      map[string]string
			err       error
		)

		vars = mux.Vars(r)
		projectId = vars["projectId"]

		proj, err = ctx.Database.Projects.FetchById(projectId)
		if err != nil {
			ctx.JSON(&rw, http.StatusNotFound, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusOK, JSON{"project": proj})
	}
}

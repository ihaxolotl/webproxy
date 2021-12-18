package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ihaxolotl/webproxy/internal/data"
)

type JSON map[string]interface{}

type CreateProjectRequest struct {
	Title       string `json:"title" validate:"required,min=6,max=64"`
	Description string `json:"description"`
}

// CreateProjectRoute is an endpoint for creating a new project.
func CreateProjectRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			req     CreateProjectRequest
			project data.Project
			err     error
		)

		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			ctx.JSON(&rw, http.StatusUnprocessableEntity, JSON{"err": err.Error()})
			return
		}

		if err = validator.New().Struct(req); err != nil {
			ctx.JSON(&rw, http.StatusUnprocessableEntity, JSON{"err": err.Error()})
			return
		}

		if err = data.InsertAndGetProject(ctx.Database, &project); err != nil {
			ctx.JSON(&rw, http.StatusInternalServerError, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusCreated, JSON{
			"msg":     "Project successfully created",
			"project": project,
		})
	}
}

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

// GetProjectsByIdRoute is an endpoint for fetching a project by its id.
func GetProjectsByIdRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			projectId string
			project   *data.Project
			vars      map[string]string
			err       error
		)

		vars = mux.Vars(r)
		projectId = vars["projectId"]

		project, err = data.GetProjectById(ctx.Database, projectId)
		if err != nil {
			ctx.JSON(&rw, http.StatusNotFound, JSON{"err": err.Error()})
			return
		}

		ctx.JSON(&rw, http.StatusOK, JSON{"project": project})
	}
}

// GetProjectsProxyRoute is an endpoint for connecting to the intercept proxy
// for the project.
func GetProjectsProxyRoute(ctx Context) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx.JSON(&rw, http.StatusNotImplemented, JSON{"msg": "Not implemented."})
	}
}

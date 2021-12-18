package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/ihaxolotl/webproxy/internal/data"
)

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

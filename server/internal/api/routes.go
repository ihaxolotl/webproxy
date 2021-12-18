package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	Database *sql.DB
}

type Route struct {
	Name    string
	URL     string
	Method  string
	Handler func(Context) http.HandlerFunc
}

func (ctx *Context) JSON(rw *http.ResponseWriter, code int, payload interface{}) {
	var (
		buf []byte
		err error
	)

	(*rw).Header().Set("Content-Type", "application/json")
	(*rw).WriteHeader(code)

	buf, err = json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	if _, err = (*rw).Write(buf); err != nil {
		log.Println(err)
	}
}

var APIRoutes []Route = []Route{
	{
		Name:    "Index",
		URL:     "/",
		Method:  "GET",
		Handler: IndexRoute,
	},
	{
		Name:    "GetProjects",
		URL:     "/projects",
		Method:  "GET",
		Handler: GetProjectsRoute,
	},
	{
		Name:    "CreateProject",
		URL:     "/projects",
		Method:  "POST",
		Handler: CreateProjectRoute,
	},
	{
		Name:    "GetProjectById",
		URL:     "/projects/{projectId}",
		Method:  "GET",
		Handler: GetProjectsByIdRoute,
	},
	{
		Name:    "GetProjectProxy",
		URL:     "/projects/{projectId}/proxy",
		Method:  "GET",
		Handler: GetProjectsProxyRoute,
	},
}

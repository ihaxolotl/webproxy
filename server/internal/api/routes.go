package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ihaxolotl/webproxy/internal/data"
)

type Route struct {
	Name    string
	URL     string
	Method  string
	Handler func(Context) http.HandlerFunc
}

type JSON map[string]interface{}

type Context struct {
	Database *data.Database
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
		Method:  http.MethodGet,
		Handler: IndexRoute,
	},
	{
		Name:    "GetProjects",
		URL:     "/projects",
		Method:  http.MethodGet,
		Handler: GetProjectsRoute,
	},
	{
		Name:    "CreateProject",
		URL:     "/projects",
		Method:  http.MethodPost,
		Handler: CreateProjectRoute,
	},
	{
		Name:    "GetProjectById",
		URL:     "/projects/{projectId}",
		Method:  http.MethodGet,
		Handler: GetProjectsByIdRoute,
	},
	{
		Name:    "GetProjectProxy",
		URL:     "/projects/{projectId}/proxy",
		Method:  http.MethodGet,
		Handler: GetProjectProxyRoute,
	},
	{
		Name:    "GetProjectHistory",
		URL:     "/projects/{projectId}/history",
		Method:  http.MethodGet,
		Handler: GetProjectHistoryRoute,
	},
	{
		Name:    "GetRequestById",
		URL:     "/requests/{requestId}",
		Method:  http.MethodGet,
		Handler: GetRequestByIdRoute,
	},
	{
		Name:    "GetResponseById",
		URL:     "/responses/{responseId}",
		Method:  http.MethodGet,
		Handler: GetResponseByIdRoute,
	},
}

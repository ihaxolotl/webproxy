package api

import (
	"fmt"
	"net/http"
)

type Route struct {
	Name    string
	URL     string
	Method  string
	Handler http.HandlerFunc
}

func IndexRoute(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Hello, world!\n")
}

var APIRoutes []Route = []Route{
	{
		Name:    "IndexRoute",
		URL:     "/",
		Method:  "GET",
		Handler: IndexRoute,
	},
}

package api

import (
	"fmt"
	"net/http"
)

const PORT = "8080"

type Handler = func(http.ResponseWriter, *http.Request) error

type route struct {
	name   string
	method map[string]Handler
}

var routeSchema = []route{
	{
		name: "/ping",
		method: map[string]Handler{
			"GET":  getPingHandler,
			"POST": postPingHandler,
		},
	},
	{
		name: "/stats",
		method: map[string]Handler{
			"GET": getStatsHandler,
		},
	},
	{
		name: "/detail",
		method: map[string]Handler{
			"GET": getDetailStatsHandler,
		},
	},
}

func (e route) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	funcHandler, ok := e.method[r.Method]

	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := funcHandler(w, r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func Start() {

	fmt.Printf("Server started on on port %s\n", PORT)

	for _, route := range routeSchema {
		http.HandleFunc(route.name, route.ServeHTTP)
	}
	http.ListenAndServe(":"+PORT, nil)

}

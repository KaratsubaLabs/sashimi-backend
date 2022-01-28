package api

import (
	"net/http"
)

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

func getPingHandler(w http.ResponseWriter, r *http.Request) error {

	w.Header().Add("content-type", "text/strings")
	w.Write([]byte("Pong!\n"))

	return nil

}

func postPingHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getStatsHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getDetailStatsHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

package api

import (
	"net/http"
)

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

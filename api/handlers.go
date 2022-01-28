package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/karatsubalabs/sashimi-backend/sashimi"
)

func getPingHandler(w http.ResponseWriter, r *http.Request) error {

	w.Header().Add("content-type", "text/strings")
	w.Write([]byte("Pong!\n"))

	return nil

}

func postPingHandler(w http.ResponseWriter, r *http.Request) error {

	var requestData postPingRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err = sashimi.Ping(requestData.Name, requestData.URL, time.Now().Unix())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	return err

}

func getStatsHandler(w http.ResponseWriter, r *http.Request) error {

	stats, err1 := sashimi.Database.GetStats()

	if err1 == nil {

		json, err2 := json.Marshal(stats)

		if err2 == nil {

			w.Header().Add("content-type", "application/json")
			w.Write(json)

		} else {

			w.WriteHeader(http.StatusInternalServerError)
			return err2

		}

	}

	w.WriteHeader(http.StatusInternalServerError)
	return err1

}

func getDetailHandler(w http.ResponseWriter, r *http.Request) error {

	var requestData getDetailRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	details, err1 := sashimi.Database.GetDetails(requestData.Name)

	if err1 == nil {

		json, err2 := json.Marshal(details)

		if err2 == nil {

			w.Header().Add("content-type", "application/json")
			w.Write(json)

		} else {

			w.WriteHeader(http.StatusInternalServerError)
			return err2

		}

	}

	w.WriteHeader(http.StatusInternalServerError)
	return err
}

package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sashayakovtseva/social-tournament-service/controller"
)

func parsePointsParam(p string) (float32, error) {
	points, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return 0, errors.New("Unable to parse points")
	}
	if points < 0 {
		return 0, errors.New("Points should be a non negative value")
	}
	return float32(points), nil
}

func GetHandler(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func PostJsonHandler(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
		} else if r.Header.Get(CONTENT_TYPE) != APPLICATION_JSON {
			http.Error(w, "", http.StatusUnsupportedMediaType)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func HandleReset(w http.ResponseWriter, r *http.Request) {
	connector, err := controller.GetConnector()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := connector.Reset(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

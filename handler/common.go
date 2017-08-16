package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sashayakovtseva/social-tournament-service/controller"
)

type HandleFuncWithErr func(http.ResponseWriter, *http.Request) error
type Middleware func(handler http.Handler) http.Handler

func parsePointsParam(p string) (float32, error) {
	points, err := strconv.ParseFloat(p, 32)
	if err != nil {
		return 0, errors.New("unable to parse points")
	}
	if points < 0 {
		return 0, errors.New("points should be a non negative value")
	}
	return float32(points), nil
}

func MethodGet(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func MethodPost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func CntTypeJson(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(CONTENT_TYPE) != APPLICATION_JSON {
			http.Error(w, "", http.StatusUnsupportedMediaType)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func HandleWithErrWrap(handler HandleFuncWithErr) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

func HandleReset(_ http.ResponseWriter, r *http.Request) error {
	return controller.Reset(r.Context())
}

func Wrap(handler HandleFuncWithErr, middleware ...Middleware) http.Handler {
	wrapped := HandleWithErrWrap(handler)
	for _, m := range middleware {
		wrapped = m(wrapped)
	}
	return wrapped
}

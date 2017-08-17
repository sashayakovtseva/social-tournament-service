package handler

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/sashayakovtseva/social-tournament-service/controller"
)

type HandleFuncWithErr func(http.ResponseWriter, *http.Request) error
type Middleware func(handler http.Handler) http.Handler
type key string

func FilterGet(handler http.Handler) http.Handler {
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

func FilterJSON(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(contentType) != applicationJSON {
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

const requestIDKey = key("reqId")

func AddRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rand.Seed(time.Now().UnixNano())
		ctx = context.WithValue(ctx, requestIDKey, rand.Int63())
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LogElapsedTime(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log(r.Context(), r.URL.String(), "elapsed time", time.Since(start))
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

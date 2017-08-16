package handler

import (
	"context"
	"fmt"
	l "log"
	"math/rand"
	"net/http"
	"time"
)

type key string

const requestIDKey = key("reqId")

// Println calls log.Printf to print to the standard logger but adds the
// request from the event context.
func log(ctx context.Context, v ...interface{}) {
	vals := make([]interface{},1, len(v)+1)
	id, ok := ctx.Value(requestIDKey).(int64)
	if !ok {
		vals[0] = "[unknown]"
	} else {
		vals[0] = fmt.Sprintf("[%20d]", id)

	}
	vals = append(vals, v...)
	l.Println(vals...)
}

func CtxWithReqId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rand.Seed(time.Now().UnixNano())
		ctx = context.WithValue(ctx, requestIDKey, rand.Int63())
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

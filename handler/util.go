package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func logWithRequertID(ctx context.Context, v ...interface{}) {
	vals := make([]interface{}, 1, len(v)+1)
	id, ok := ctx.Value(requestIDKey).(int64)
	if !ok {
		vals[0] = "[unknown]"
	} else {
		vals[0] = fmt.Sprintf("[%20d]", id)

	}
	vals = append(vals, v...)
	log.Println(vals...)
}

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

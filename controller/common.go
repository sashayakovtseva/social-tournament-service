package controller

import (
	"context"

	"github.com/sashayakovtseva/social-tournament-service/database"
)

func Reset(ctx context.Context) (e error) {
	return database.Reset()
}

func Close() {
	if pCSingleton != nil {
		pCSingleton.Close()
	}
	if tCSingleton != nil {
		tCSingleton.Close()
	}
}

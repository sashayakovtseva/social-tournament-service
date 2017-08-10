package controller

import "github.com/sashayakovtseva/social-tournament-service/database"

func Reset() (e error) {
	return database.Reset()
}

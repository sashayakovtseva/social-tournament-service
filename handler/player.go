package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/controller"
)

const (
	DEPLOY_PORT      = 8080
	PLAYER_ID_PARAM  = "playerId"
	POINTS_PARAM     = "points"
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE     = "Content-Type"
)

func HandleTake(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()

	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err == nil {
		// todo mb create global controller?
		var playerController *controller.PlayerController
		if playerController, err = controller.GetPlayerController(); err == nil {
			err = playerController.Take(playerId, points)
		}
	}
}

func HandleFund(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()

	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err == nil {
		var playerController *controller.PlayerController
		if playerController, err = controller.GetPlayerController(); err == nil {
			err = playerController.Fund(playerId, points)
		}
	}
}

func HandleBalance(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	playerController, err := controller.GetPlayerController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// todo separate response struct?
	player := playerController.GetById(playerId)
	if player == nil {
		http.Error(w, "No such player", http.StatusBadRequest)
	} else {
		w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
		json.NewEncoder(w).Encode(player)
	}
}

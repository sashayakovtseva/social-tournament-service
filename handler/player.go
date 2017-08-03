package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/controller"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

const (
	DEPLOY_PORT      = 8080
	PLAYER_ID_PARAM  = "playerId"
	POINTS_PARAM     = "points"
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE     = "Content-Type"
)

func HandleTake(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerController, err := controller.GetPlayerController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player := playerController.GetById(playerId)
	if player == nil {
		http.Error(w, "No such player", http.StatusBadRequest)
		return
	}
	if player.Take(points) {
		playerController.Update(player)
	} else {
		http.Error(w, "Not enough points", http.StatusBadRequest)
	}
}

func HandleFund(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerController, err := controller.GetPlayerController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player := playerController.GetById(playerId)
	if player == nil {
		playerController.Create(entity.NewPlayer(playerId, points))
	} else {
		player.Fund(points)
		playerController.Update(player)
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
	player := playerController.GetById(playerId)
	if player == nil {
		http.Error(w, "No such player", http.StatusBadRequest)
	} else {
		w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
		json.NewEncoder(w).Encode(player)
	}
}

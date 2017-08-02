package view

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sashayakovtseva/social-tournament-service/controller"
	"github.com/sashayakovtseva/social-tournament-service/model"
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
	points, err := strconv.Atoi(params.Get(POINTS_PARAM))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if points < 0 {
		http.Error(w, "Points should be a non negative value", http.StatusBadRequest)
		return
	}
	playerController, err := controller.GetPlayerController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		player := playerController.GetPlayerById(playerId)
		if player == nil {
			playerController.InsertPlayer(model.NewPlayer(playerId, points))
		} else {
			player.Fund(points)
			playerController.UpdatePlayer(player)
		}
	}
}

func HandleFund(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := strconv.Atoi(params.Get(POINTS_PARAM))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if points < 0 {
		http.Error(w, "Points should be a non negative value", http.StatusBadRequest)
		return
	}
	playerController, err := controller.GetPlayerController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		player := playerController.GetPlayerById(playerId)
		if player == nil {
			playerController.InsertPlayer(model.NewPlayer(playerId, points))
		} else {
			player.Fund(points)
			playerController.UpdatePlayer(player)
		}
	}
}
func HandleBalance(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	playerController, err := controller.GetPlayerController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		player := playerController.GetPlayerById(playerId)
		if player == nil {
			http.Error(w, "No such player", http.StatusBadRequest)
		} else {
			w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
			json.NewEncoder(w).Encode(player)
		}
	}
}

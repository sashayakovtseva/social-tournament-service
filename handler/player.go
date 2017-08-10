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

type PlayerBalanceResponse struct {
	Id      string  `json:"playerId"`
	Balance float32 `json:"balance"`
}

func HandleTake(_ http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err != nil {
		return err
	}
	playerController := controller.GetPlayerController()
	return playerController.Take(playerId, points)
}

func HandleFund(_ http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err != nil {
		return err
	}
	playerController := controller.GetPlayerController()
	return playerController.Fund(playerId, points)
}

func HandleBalance(w http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	playerController := controller.GetPlayerController()
	var balance float32
	balance, err := playerController.Balance(playerId)
	if err != nil {
		return err
	}
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	return json.NewEncoder(w).Encode(PlayerBalanceResponse{Balance: balance, Id: playerId})
}

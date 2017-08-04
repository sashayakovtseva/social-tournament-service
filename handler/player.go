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
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()

	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	var playerController *controller.PlayerController
	if playerController, err = controller.GetPlayerController(); err == nil {
		var result *entity.PlayerBalanceResponse
		result, err = playerController.Balance(playerId)
		if err == nil {
			w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
			err = json.NewEncoder(w).Encode(result)
		}
	}
}

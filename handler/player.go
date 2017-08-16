package handler

import (
	"encoding/json"
	"errors"
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

func HandleTake(w http.ResponseWriter, r *http.Request) error {
	return handleUpdate(w, r, controller.GetPlayerController().Take)
}

func HandleFund(w http.ResponseWriter, r *http.Request) error {
	return handleUpdate(w, r, controller.GetPlayerController().Fund)
}

func HandleBalance(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	playerController := controller.GetPlayerController()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log(ctx, err.Error())
		return err
	case player := <-playerController.Read(playerId):
		if player != nil {
			w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
			return json.NewEncoder(w).Encode(PlayerBalanceResponse{Balance: player.Balance(), Id: player.Id()})
		}
		return errors.New("no such player")
	}
}

func handleUpdate(_ http.ResponseWriter, r *http.Request, update func(string, float32) chan error) error {
	ctx := r.Context()
	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		err := ctx.Err()
		log(ctx, err.Error())
		return err
	case err := <-update(playerId, points):
		return err
	}
}

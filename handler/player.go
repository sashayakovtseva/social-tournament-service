package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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
	ctx := r.Context()
	log(ctx, "take started")
	defer log(ctx, "take ended")

	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err != nil {
		return err
	}
	playerController := controller.GetPlayerController()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log(ctx, err.Error())
		return err
	case err := <-playerController.Take(playerId, points):
		return err
	}
}

func HandleFund(_ http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	log(ctx, "fund started")
	defer log(ctx, "fund ended")

	params := r.URL.Query()
	playerId := params.Get(PLAYER_ID_PARAM)
	points, err := parsePointsParam(params.Get(POINTS_PARAM))
	if err != nil {
		return err
	}
	playerController := controller.GetPlayerController()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log(ctx, err.Error())
		return err
	case err := <-playerController.Fund(playerId, points):
		return err
	}
}

func HandleBalance(w http.ResponseWriter, r *http.Request) error {
	start := time.Now()
	ctx := r.Context()
	log(ctx, "balance started")
	defer func() {
		log(ctx, "balance ended, time elapsed", time.Since(start))
	}()
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

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/controller"
)

const (
	DeployPort      = 8080
	playerIDParam   = "playerId"
	pointsParam     = "points"
	applicationJSON = "application/json"
	contentType     = "Content-Type"
)

type PlayerBalanceResponse struct {
	ID      string  `json:"playerId"`
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
	playerID := params.Get(playerIDParam)
	playerController := controller.GetPlayerController()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log(ctx, err.Error())
		return err
	case player := <-playerController.Read(playerID):
		if player != nil {
			w.Header().Set(contentType, applicationJSON)
			return json.NewEncoder(w).Encode(PlayerBalanceResponse{Balance: player.Balance(), ID: player.ID()})
		}
		return errors.New("no such player")
	}
}

func handleUpdate(_ http.ResponseWriter, r *http.Request, update func(string, float32) chan error) error {
	ctx := r.Context()
	params := r.URL.Query()
	playerID := params.Get(playerIDParam)
	points, err := parsePointsParam(params.Get(pointsParam))
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		err := ctx.Err()
		log(ctx, err.Error())
		return err
	case err := <-update(playerID, points):
		return err
	}
}

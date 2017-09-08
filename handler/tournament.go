package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/controller"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

const (
	tournamentIDParam = "tournamentId"
	depositParam      = "deposit"
	backerIDParam     = "backerId"
)

func HandleAnnounce(_ http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := r.URL.Query()
	tournamentID := params.Get(tournamentIDParam)
	deposit, err := parsePointsParam(params.Get(depositParam))
	if err != nil {
		return err
	}
	tournamentController := controller.GetTournamentController()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		logWithRequertID(ctx, err.Error())
		return err
	case err := <-tournamentController.Announce(tournamentID, deposit):
		return err
	}
}

func HandleJoin(_ http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	tournamentID := params.Get(tournamentIDParam)
	playerID := params.Get(playerIDParam)
	backers := params[backerIDParam]
	tournamentController := controller.GetTournamentController()
	return tournamentController.JoinTournament(r.Context(), tournamentID, playerID, backers)
}

func HandleResult(_ http.ResponseWriter, r *http.Request) error {
	var requestBody *entity.ResultTournamentRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return err
	}
	tournamentController := controller.GetTournamentController()
	return tournamentController.Result(r.Context(), requestBody)
}

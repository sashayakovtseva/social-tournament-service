package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/controller"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

const (
	TOURNAMENT_ID_PARAM = "tournamentId"
	DEPOSIT_PARAM       = "deposit"
	BACKER_ID_PARAM     = "backerId"
)

func HandleAnnounce(_ http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := r.URL.Query()
	tournamentId := params.Get(TOURNAMENT_ID_PARAM)
	deposit, err := parsePointsParam(params.Get(DEPOSIT_PARAM))
	if err != nil {
		return err
	}
	tournamentController := controller.GetTournamentController()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log(ctx, err.Error())
		return err
	case err := <-tournamentController.Announce(tournamentId, deposit):
		return err
	}
}

func HandleJoin(_ http.ResponseWriter, r *http.Request) error {
	params := r.URL.Query()
	tournamentId := params.Get(TOURNAMENT_ID_PARAM)
	playerId := params.Get(PLAYER_ID_PARAM)
	backers := params[BACKER_ID_PARAM]
	tournamentController := controller.GetTournamentController()
	return tournamentController.JoinTournament(r.Context(), tournamentId, playerId, backers)
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

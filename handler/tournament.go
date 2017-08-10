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

func HandleAnnounce(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()

	params := r.URL.Query()
	tournamentId := params.Get(TOURNAMENT_ID_PARAM)
	deposit, err := parsePointsParam(params.Get(DEPOSIT_PARAM))
	if err == nil {
		tournamentController := controller.GetTournamentController()
		err = tournamentController.Announce(tournamentId, deposit)
	}
}

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()

	params := r.URL.Query()
	tournamentId := params.Get(TOURNAMENT_ID_PARAM)
	playerId := params.Get(PLAYER_ID_PARAM)
	backers := params[BACKER_ID_PARAM]

	tournamentController := controller.GetTournamentController()
	err = tournamentController.JoinTournament(tournamentId, playerId, backers)
}

func HandleResult(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}()

	var requestBody *entity.ResultTournamentRequest
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err == nil {
		tournamentController := controller.GetTournamentController()
		err = tournamentController.Result(requestBody)
	}
}

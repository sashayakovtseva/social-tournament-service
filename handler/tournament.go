package handler

import (
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
	params := r.URL.Query()
	tournamentId := params.Get(TOURNAMENT_ID_PARAM)
	deposit, err := parsePointsParam(params.Get(DEPOSIT_PARAM))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tournamentController, err := controller.GetTournamentController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournament := tournamentController.GetById(tournamentId)
	if tournament != nil {
		http.Error(w, "Tournament already exists", http.StatusBadRequest)
		return
	}
	tournamentController.Create(entity.NewTournament(tournamentId, deposit))
}

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	tournamentId := params.Get(TOURNAMENT_ID_PARAM)
	playerId := params.Get(PLAYER_ID_PARAM)
	backers := params[BACKER_ID_PARAM]

	tournamentController, err := controller.GetTournamentController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerController, err := controller.GetPlayerController()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tournament := tournamentController.GetById(tournamentId)
	if tournament == nil {
		http.Error(w, "No such tournament", http.StatusBadRequest)
		return
	}

	player := playerController.GetById(playerId)
	if player == nil {
		http.Error(w, "No such player", http.StatusBadRequest)
		return
	}

	// todo do we really need this? db constraint should work
	backPlayers := make([]*entity.Player, 0)
	for _, backerId := range backers {
		backer := playerController.GetById(backerId)
		if backer == nil {
			http.Error(w, "One or more backers are not found", http.StatusBadRequest)
			return
		}
		backPlayers = append(backPlayers, backer)
	}

	err = tournamentController.JoinTournament(tournament, player, backPlayers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

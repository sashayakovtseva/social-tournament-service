package controller

import (
	"context"
	"errors"
	"sync"

	db "github.com/sashayakovtseva/social-tournament-service/database"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	tournControllerSingleton *TournamentController
	tournControllerOnce      sync.Once
)

type TournamentController struct {
	lock sync.Mutex
}

func GetTournamentController() *TournamentController {
	tournControllerOnce.Do(func() {
		tournControllerSingleton = new(TournamentController)
	})
	return tournControllerSingleton
}

func (tC *TournamentController) Announce(ctx context.Context, tournamentId string, deposit float32) error {
	tC.lock.Lock()
	defer tC.lock.Unlock()

	tournament := db.TournamentConn.Read(tournamentId)
	if tournament != nil {
		return errors.New("tournament already exists")
	}
	return db.TournamentConn.Create(entity.NewTournament(tournamentId, deposit, false))
}

func (tC *TournamentController) Result(ctx context.Context, tournamentResult *entity.ResultTournamentRequest) error {
	tC.lock.Lock()
	defer tC.lock.Unlock()

	tournament := db.TournamentConn.Read(tournamentResult.Id)
	if tournament == nil {
		return errors.New("no such tournament")
	}
	if tournament.IsFinished() {
		return errors.New("tournament has finished already")
	}
	participants, backPlayers, err := db.TournamentConn.SelectParticipants(tournamentResult.Id)
	if err != nil {
		return err
	}
	participantsSet := make(map[string]float32, len(participants))
	for _, p := range participants {
		participantsSet[p.Id()] = p.Balance()
	}
	winners := make(map[entity.Player]float32)
	for _, winner := range tournamentResult.Winners {
		if balance, ok := participantsSet[winner.Id]; !ok {
			return errors.New("winner is not a participant")
		} else {
			winners[*entity.NewPlayer(winner.Id, balance)] = winner.Prize
		}
	}
	involved := tournament.Result(participants, backPlayers, winners)
	return db.TournamentConn.UpdateResult(tournament, involved)
}

func (tC *TournamentController) JoinTournament(ctx context.Context, tournamentId, playerId string, backersId []string) error {
	tC.lock.Lock()
	defer tC.lock.Unlock()

	tournament := db.TournamentConn.Read(tournamentId)
	if tournament == nil {
		return errors.New("no such tournament")
	}
	if tournament.IsFinished() {
		return errors.New("tournament has finished already")
	}
	player := db.PlayerConn.Read(playerId)
	if player == nil {
		return errors.New("no such player")
	}
	backPlayers := make([]*entity.Player, 0)
	for _, backerId := range backersId {
		backer := db.PlayerConn.Read(backerId)
		if backer == nil {
			return errors.New("one or more backers are not found")
		}
		backPlayers = append(backPlayers, backer)
	}
	if !tournament.Join(player, backPlayers) {
		return errors.New("not enough points to join")
	}
	return db.TournamentConn.UpdateJoin(tournament, player, backPlayers)
}

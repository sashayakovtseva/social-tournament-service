package controller

import (
	"context"
	"errors"
	"sync"

	db "github.com/sashayakovtseva/social-tournament-service/database"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	tCSingleton *TournamentController
	tCOnce      sync.Once
)

type TournamentController struct {
	lock sync.Mutex
}

func GetTournamentController() *TournamentController {
	tCOnce.Do(func() {
		tCSingleton = new(TournamentController)
	})
	return tCSingleton
}

func (tC *TournamentController) Announce(tournamentID string, deposit float32) chan error {
	err := make(chan error, 1)
	go func() {
		err <- db.TournamentConn.Create(entity.NewTournament(tournamentID, deposit, false))
	}()
	return err
}

func (tC *TournamentController) Close() {}

func (tC *TournamentController) Result(ctx context.Context, tournamentResult *entity.ResultTournamentRequest) error {
	tC.lock.Lock()
	defer tC.lock.Unlock()

	tournament := db.TournamentConn.Read(tournamentResult.ID)
	if tournament == nil {
		return errors.New("no such tournament")
	}
	if tournament.IsFinished() {
		return errors.New("tournament has finished already")
	}
	participants, backPlayers, err := db.TournamentConn.SelectParticipants(tournamentResult.ID)
	if err != nil {
		return err
	}
	participantsSet := make(map[string]float32, len(participants))
	for _, p := range participants {
		participantsSet[p.ID()] = p.Balance()
	}
	winners := make(map[entity.Player]float32)
	for _, winner := range tournamentResult.Winners {
		balance, ok := participantsSet[winner.ID]
		if !ok {
			return errors.New("winner is not a participant")
		}
		winners[*entity.NewPlayer(winner.ID, balance)] = winner.Prize
	}
	involved := tournament.Result(participants, backPlayers, winners)
	return db.TournamentConn.UpdateResult(tournament, involved)
}

func (tC *TournamentController) JoinTournament(ctx context.Context, tournamentID, playerID string, backersID []string) error {
	tC.lock.Lock()
	defer tC.lock.Unlock()

	tournament := db.TournamentConn.Read(tournamentID)
	if tournament == nil {
		return errors.New("no such tournament")
	}
	if tournament.IsFinished() {
		return errors.New("tournament has finished already")
	}
	player := db.PlayerConn.Read(playerID)
	if player == nil {
		return errors.New("no such player")
	}
	backPlayers := make([]*entity.Player, 0)
	for _, backerId := range backersID {
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

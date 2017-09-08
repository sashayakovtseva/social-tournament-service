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
	sync.Mutex
}


func GetTournamentController() *TournamentController {
	tCOnce.Do(func() {
		tCSingleton = new(TournamentController)
	})
	return tCSingleton
}

func (controller *TournamentController) Announce(tournamentID string, deposit float32) chan error {
	err := make(chan error, 1)
	go func() {
		err <- db.TournamentConn.Create(entity.NewTournament(tournamentID, deposit, false))
	}()
	return err
}

func (controller *TournamentController) Close() {}

func (controller *TournamentController) Result(ctx context.Context, tournamentResult *entity.ResultTournamentRequest) error {
	controller.Lock()
	defer controller.Unlock()

	tournament := db.TournamentConn.Read(tournamentResult.ID)
	if tournament == nil {
		return errors.New("no such tournament")
	}
	if tournament.IsFinished() {
		return errors.New("tournament has finished already")
	}

	participants, err := db.TournamentConn.SelectParticipants(tournamentResult.ID)
	if err != nil {
		return err
	}

	backPlayers, err := controller.selectAllBackPlayers(tournamentResult.ID, participants)
	if err != nil {
		return err
	}

	participantsSet := make(map[string]float32, len(participants))
	for _, p := range participants {
		participantsSet[p.ID()] = p.Balance()
	}

	winners, err := controller.checkWinnersAndFormMap(participantsSet, tournamentResult.Winners)
	if err != nil {
		return err
	}

	involved := tournament.Result(participants, backPlayers, winners)
	return db.TournamentConn.UpdateResult(tournament, involved)
}

func (controller *TournamentController) selectAllBackPlayers(tournamentID string, participants []*entity.Player) ([][]*entity.Player, error) {
	backPlayers := make([][]*entity.Player, 0, len(participants))
	for _, p := range participants {
		backs, err := db.TournamentConn.SelectBackPlayers(tournamentID, p.ID())
		if err != nil {
			return nil, err
		}
		backPlayers = append(backPlayers, backs)
	}
	return backPlayers, nil
}

func (controller *TournamentController) checkWinnersAndFormMap(participantsSet map[string]float32, winners []entity.Winner) (map[entity.Player]float32, error) {
	w := make(map[entity.Player]float32)
	for _, winner := range winners {
		balance, ok := participantsSet[winner.ID]
		if !ok {
			return nil, errors.New("winner is not a participant")
		}
		w[*entity.NewPlayer(winner.ID, balance)] = winner.Prize
	}
	return w, nil
}

func (controller *TournamentController) JoinTournament(ctx context.Context, tournamentID, playerID string, backersID []string) error {
	controller.Lock()
	defer controller.Unlock()

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
	backPlayers, err := controller.readBackPlayersByIDs(backersID)
	if err != nil {
		return err
	}
	if !tournament.Join(player, backPlayers) {
		return errors.New("not enough points to join")
	}
	return db.TournamentConn.UpdateJoin(tournament, player, backPlayers)
}

func (controller *TournamentController) readBackPlayersByIDs(backersID []string) ([]*entity.Player, error) {
	backPlayers := make([]*entity.Player, 0)
	for _, backerId := range backersID {
		backer := db.PlayerConn.Read(backerId)
		if backer == nil {
			return nil, errors.New("one or more backers are not found")
		}
		backPlayers = append(backPlayers, backer)
	}
	return backPlayers, nil
}

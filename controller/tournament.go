package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	tournControllerSingleton *TournamentController
	tournControllerError     error
	tournControllerOnce      sync.Once
)

type TournamentController struct {
	connector                  *DBConnector
	preparedInsertTournament   *sql.Stmt
	preparedUpdateTournament   *sql.Stmt
	preparedSelectTournament   *sql.Stmt
	preparedJoinTournament     *sql.Stmt
	preparedInsertBackPlayer   *sql.Stmt
	preparedSelectParticipants *sql.Stmt
	preparedSelectBackers      *sql.Stmt
}

func newTournamentController() (*TournamentController, error) {
	connector, err := GetConnector()
	if err != nil {
		return nil, err
	}

	tournController := &TournamentController{connector: connector}
	tournController.preparedInsertTournament, err = connector.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		TOURNAMENTS_TABLE_NAME, TOURNAMENT_ID_COL_NAME, DEPOSIT_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedUpdateTournament, err = connector.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		TOURNAMENTS_TABLE_NAME, FINISHED_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedSelectTournament, err = connector.Prepare(fmt.Sprintf(
		`SELECT %s, %s, %s FROM %s WHERE %s = ?`,
		TOURNAMENT_ID_COL_NAME, DEPOSIT_COL_NAME, FINISHED_COL_NAME, TOURNAMENTS_TABLE_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedJoinTournament, err = connector.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		P2T_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedInsertBackPlayer, err = connector.Prepare(fmt.Sprintf(
		`INSERT INTO %s(%s,%s, %s) values(?,?, ?)`,
		P2B_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME, BACKER_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedSelectParticipants, err = connector.Prepare(fmt.Sprintf(
		`SELECT %s, %s FROM %s INNER JOIN %s  USING (%s) WHERE %s = ?`,
		PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, P2T_TABLE_NAME,
		PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedSelectBackers, err = connector.Prepare(fmt.Sprintf(
		`SELECT %s, %s FROM %s INNER JOIN %s ON %s.%s  = %s.%s WHERE %s = ? AND %s.%s = ?`,
		PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, P2B_TABLE_NAME, P2B_TABLE_NAME, BACKER_COL_NAME,
		PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME, P2B_TABLE_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	return tournController, nil
}

func GetTournamentController() (*TournamentController, error) {
	tournControllerOnce.Do(func() {
		tournControllerSingleton, tournControllerError = newTournamentController()
	})
	return tournControllerSingleton, tournControllerError
}

func (tC *TournamentController) Close() {
	if tC.preparedInsertTournament != nil {
		tC.preparedInsertTournament.Close()
	}
	if tC.preparedUpdateTournament != nil {
		tC.preparedUpdateTournament.Close()
	}
	if tC.preparedSelectTournament != nil {
		tC.preparedSelectTournament.Close()
	}
	if tC.preparedJoinTournament != nil {
		tC.preparedJoinTournament.Close()
	}

	if tC.preparedInsertBackPlayer != nil {
		tC.preparedInsertBackPlayer.Close()
	}
	if tC.preparedSelectParticipants != nil {
		tC.preparedSelectParticipants.Close()
	}
	if tC.preparedSelectBackers != nil {
		tC.preparedSelectBackers.Close()
	}
}

func (tC *TournamentController) GetById(id string) *entity.Tournament {
	var tournamentId string
	var tournamentDeposit float32
	var tournamentIsFinished bool
	err := tC.preparedSelectTournament.QueryRow(id).Scan(&tournamentId, &tournamentDeposit, &tournamentIsFinished)
	if err != nil {
		return nil
	}
	return entity.NewTournament(tournamentId, tournamentDeposit, tournamentIsFinished)
}

func (tC *TournamentController) Update(tournament *entity.Tournament) error {
	_, err := tC.preparedUpdateTournament.Exec(tournament.IsFinished(), tournament.Id())
	return err
}

func (tC *TournamentController) Create(tournament *entity.Tournament) error {
	_, err := tC.preparedInsertTournament.Exec(tournament.Id(), tournament.Deposit())
	return err
}

func (tC *TournamentController) Announce(tournamentId string, deposit float32) error {
	tournament := tC.GetById(tournamentId)
	if tournament != nil {
		return errors.New("Tournament already exists")
	}
	tC.Create(entity.NewTournament(tournamentId, deposit, false))
	return nil
}

func (tC *TournamentController) Result(tournamentResult *entity.ResultTournamentRequest) error {
	tournament := tC.GetById(tournamentResult.Id)
	if tournament == nil {
		return errors.New("No such tournament")
	}
	if tournament.IsFinished() {
		return errors.New("Tournament has finished already")
	}

	playerController, err := GetPlayerController()
	if err != nil {
		return err
	}

	participants := make([]*entity.Player, 0)
	participantsSet := make(map[string]float32)
	backPlayers := make([][]*entity.Player, 0)

	rows, err := tC.preparedSelectParticipants.Query(tournamentResult.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var playerId string
	var balance float32
	for rows.Next() {
		rows.Scan(&playerId, &balance)
		participantsSet[playerId] = balance
		participants = append(participants, entity.NewPlayer(playerId, balance))

		r, err := tC.preparedSelectBackers.Query(tournamentResult.Id, playerId)
		if err != nil {
			r.Close()
			return err
		}

		var backers []*entity.Player
		for r.Next() {
			r.Scan(&playerId, &balance)
			backers = append(backers, entity.NewPlayer(playerId, balance))
		}
		r.Close()

		backPlayers = append(backPlayers, backers)
	}

	winners := make(map[entity.Player]float32)
	for _, winner := range tournamentResult.Winners {
		if balance, ok := participantsSet[winner.Id]; !ok {
			return errors.New("Winner is not a participant")
		} else {
			winners[*entity.NewPlayer(winner.Id, balance)] = winner.Prize
		}
	}

	involved := tournament.Result(participants, backPlayers, winners)
	playerController.Update(involved...)
	tC.Update(tournament)
	return nil
}

func (tC *TournamentController) JoinTournament(tournamentId, playerId string, backersId []string) error {
	tournament := tC.GetById(tournamentId)
	if tournament == nil {
		return errors.New("No such tournament")
	}
	if tournament.IsFinished() {
		return errors.New("Tournament has finished already")
	}
	playerController, err := GetPlayerController()
	if err != nil {
		return err
	}
	player := playerController.GetById(playerId)
	if player == nil {
		return errors.New("No such player")
	}
	// todo do we really need this? db constraint should work
	backPlayers := make([]*entity.Player, 0)
	for _, backerId := range backersId {
		backer := playerController.GetById(backerId)
		if backer == nil {
			return errors.New("One or more backers are not found")
		}
		backPlayers = append(backPlayers, backer)
	}

	if !tournament.Join(player, backPlayers) {
		return errors.New("Not enough points to join")
	}

	tx, err := tC.connector.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Stmt(tC.preparedJoinTournament).Exec(player.Id(), tournament.Id())
	if err != nil {
		tx.Rollback()
		return err
	}
	txPreparedBackPlayer := tx.Stmt(tC.preparedInsertBackPlayer)
	for _, backer := range backPlayers {
		_, err := txPreparedBackPlayer.Exec(player.Id(), tournament.Id(), backer.Id())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

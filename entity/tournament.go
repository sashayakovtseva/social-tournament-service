package entity

import "database/sql"

type Tournament struct {
	Id       string
	Deposit  int
	WinnerId sql.NullString
}

func NewTournament(id string, deposit int) *Tournament {
	return &Tournament{
		Id:      id,
		Deposit: deposit,
	}
}

func (t *Tournament) Result() {

}

func (t *Tournament) Join(player *Player, backers []*Player) bool {
	totalBalance := player.Balance
	for _, backer := range backers {
		totalBalance += backer.Balance
	}

	return totalBalance >= t.Deposit
}

package model

type Tournament struct {
	id      string
	deposit int
	winner  string
}

func NewTournament(id string, deposit int) *Tournament {
	return &Tournament{
		id:      id,
		deposit: deposit,
	}
}

func (t *Tournament) Result() {

}
package model

type Player struct {
	Id      string `json:"playerId"`
	Balance int    `json:"balance"`
}

func NewPlayer(id string, balance int) *Player {
	return &Player{
		Id:      id,
		Balance: balance,
	}
}

func (p *Player) Take(points int) bool {
	if points > p.Balance {
		return false
	} else {
		p.Balance -= points
		return true
	}
}

func (p *Player) Fund(points int) bool {
	p.Balance += points
	return true
}

package entity

type Player struct {
	id      string
	balance float32
}

func NewPlayer(id string, balance float32) *Player {
	return &Player{
		id:      id,
		balance: balance,
	}
}

func (p *Player) Take(points float32) bool {
	if points > p.balance {
		return false
	} else {
		p.balance -= points
		return true
	}
}

func (p *Player) Fund(points float32) {
	p.balance += points
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Balance() float32 {
	return p.balance
}

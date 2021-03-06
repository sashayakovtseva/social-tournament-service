package entity

type ResultTournamentRequest struct {
	Id      string `json:"tournamentId"`
	Winners []struct {
		Id    string  `json:"playerId"`
		Prize float32 `json:"prize"`
	} `json:"winners"`
}

type Tournament struct {
	id         string
	deposit    float32
	isFinished bool
}

func NewTournament(id string, deposit float32, isFinished bool) *Tournament {
	return &Tournament{
		id:         id,
		deposit:    deposit,
		isFinished: isFinished,
	}
}

func (t *Tournament) Result(participants []*Player, backPlayers [][]*Player, winners map[Player]float32) []*Player {
	involved := make(map[string]float32)
	for i, player := range participants {
		backers := backPlayers[i]
		if prize, ok := winners[*player]; ok {
			won := prize / float32(len(backers)+1)
			involved[player.Id()] = player.Balance() + won
			for _, back := range backers {
				involved[back.Id()] = back.Balance() + won
			}
		}
	}
	t.isFinished = true

	result := make([]*Player, 0, len(involved))
	for id, balance := range involved {
		result = append(result, NewPlayer(id, balance))
	}

	return result
}

func (t *Tournament) Join(player *Player, backers []*Player) bool {
	contribute := t.deposit / float32(len(backers)+1)
	if !player.Take(contribute) {
		return false
	}
	for _, backer := range backers {
		if !backer.Take(contribute) {
			return false
		}
	}
	return true
}

func (t *Tournament) Id() string {
	return t.id
}
func (t *Tournament) Deposit() float32 {
	return t.deposit
}
func (t *Tournament) IsFinished() bool {
	return t.isFinished
}

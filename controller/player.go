package controller

import (
	"sync"

	db "github.com/sashayakovtseva/social-tournament-service/database"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	pCSingleton *PlayerController
	pCOnce      sync.Once
)

const (
	OP_TAKE = iota
	OP_FUND
)

type playerJob struct {
	playerId string
	op       int
	points   float32
}

type PlayerController struct {
	done       chan struct{}
	jobs       chan playerJob
	jobResults chan error
}

func GetPlayerController() *PlayerController {
	pCOnce.Do(func() {
		pCSingleton = new(PlayerController)
		pCSingleton.done = make(chan struct{})
		pCSingleton.jobs = make(chan playerJob)
		pCSingleton.jobResults = make(chan error)

		go pCSingleton.listenUpdate()
	})
	return pCSingleton
}

func (pC *PlayerController) Read(playerId string) chan *entity.Player {
	player := make(chan *entity.Player, 1)
	go func() {
		player <- db.PlayerConn.Read(playerId)
	}()
	return player
}

func (pC *PlayerController) Fund(playerId string, points float32) chan error {
	err := make(chan error, 1)
	go func() {
		pC.jobs <- playerJob{playerId, OP_FUND, points}
		e := <-pC.jobResults
		if e == db.ErrNoSuchPlayer {
			err <- db.PlayerConn.Create(entity.NewPlayer(playerId, points))
			return
		}
	}()
	return err
}

func (pC *PlayerController) Take(playerId string, points float32) chan error {
	err := make(chan error, 1)
	go func() {
		pC.jobs <- playerJob{playerId, OP_TAKE, points}
		err <- <-pC.jobResults
	}()
	return err
}

func (pC *PlayerController) Close() {
	close(pC.done)
}

// only one goroutine executes take & fund
// no locks
func (pC *PlayerController) listenUpdate() {
	for {
		select {
		case <-pC.done:
			return
		case job := <-pC.jobs:
			switch job.op {
			case OP_TAKE:
				pC.jobResults <- db.PlayerConn.Take(job.playerId, job.points)
			case OP_FUND:
				pC.jobResults <- db.PlayerConn.Fund(job.playerId, job.points)
			}
		}
	}
}

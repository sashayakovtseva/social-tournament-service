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
	opTake = iota
	opFund
)

type playerJob struct {
	playerID string
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

func (controller *PlayerController) Read(playerID string) chan *entity.Player {
	player := make(chan *entity.Player, 1)
	go func() {
		player <- db.PlayerConn.Read(playerID)
	}()
	return player
}

func (controller *PlayerController) FundPoints(playerID string, points float32) chan error {
	err := make(chan error, 1)
	go func() {
		controller.jobs <- playerJob{playerID, opFund, points}
		e := <-controller.jobResults
		if e == db.ErrNoSuchPlayer {
			err <- db.PlayerConn.Create(entity.NewPlayer(playerID, points))
			return
		}
	}()
	return err
}

func (controller *PlayerController) TakePoints(playerID string, points float32) chan error {
	err := make(chan error, 1)
	go func() {
		controller.jobs <- playerJob{playerID, opTake, points}
		err <- <-controller.jobResults
	}()
	return err
}

func (controller *PlayerController) Close() {
	close(controller.done)
}

// only one goroutine executes take & fund
// no locks
func (controller *PlayerController) listenUpdate() {
	for {
		select {
		case <-controller.done:
			return
		case job := <-controller.jobs:
			switch job.op {
			case opTake:
				controller.jobResults <- db.PlayerConn.TakePoints(job.playerID, job.points)
			case opFund:
				controller.jobResults <- db.PlayerConn.FundPoints(job.playerID, job.points)
			}
		}
	}
}

package controller

import (
	"context"
	"errors"
	"sync"

	db "github.com/sashayakovtseva/social-tournament-service/database"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	playerControllerSingleton *PlayerController
	playerControllerOnce      sync.Once
)

type PlayerController struct {
	lock sync.Mutex
}

func GetPlayerController() *PlayerController {
	playerControllerOnce.Do(func() {
		playerControllerSingleton = new(PlayerController)
	})
	return playerControllerSingleton
}

func (pC *PlayerController) Take(ctx context.Context, playerId string, points float32) error {
	pC.lock.Lock()
	defer pC.lock.Unlock()

	player := db.PlayerConn.Read(playerId)
	if player == nil {
		return errors.New("no such player")
	}
	if player.Take(points) {
		db.PlayerConn.Update(player)
		return nil
	}
	return errors.New("not enough points")
}

func (pC *PlayerController) Fund(ctx context.Context, playerId string, points float32) error {
	pC.lock.Lock()
	defer pC.lock.Unlock()

	player := db.PlayerConn.Read(playerId)
	if player == nil {
		return db.PlayerConn.Create(entity.NewPlayer(playerId, points))
	}
	player.Fund(points)
	return db.PlayerConn.Update(player)
}

func (pC *PlayerController) Balance(ctx context.Context, playerId string) (float32, error) {
	player := db.PlayerConn.Read(playerId)
	if player == nil {
		return 0, errors.New("no such player")
	}
	return player.Balance(), nil
}

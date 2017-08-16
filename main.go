package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sashayakovtseva/social-tournament-service/controller"
	"github.com/sashayakovtseva/social-tournament-service/database"
	h "github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {

	http.Handle("/take", h.Wrap(h.HandleTake,
		h.LogElapsedTime,
		h.AddRequestId,
		h.FilterGet,
	))
	http.Handle("/fund", h.Wrap(h.HandleFund,
		h.LogElapsedTime,
		h.AddRequestId,
		h.FilterGet,
	))
	http.Handle("/balance", h.Wrap(h.HandleBalance,
		h.LogElapsedTime,
		h.AddRequestId,
		h.FilterGet,
	))
	http.Handle("/announceTournament", h.Wrap(h.HandleAnnounce,
		h.LogElapsedTime,
		h.AddRequestId,
		h.FilterGet,
	))
	http.Handle("/joinTournament", h.Wrap(h.HandleJoin,
		h.LogElapsedTime,
		h.AddRequestId,
		h.FilterGet,
	))
	http.Handle("/resultTournament", h.Wrap(h.HandleResult,
		h.LogElapsedTime,
		h.AddRequestId,
		h.FilterJson,
		h.MethodPost,
	))
	http.Handle("/reset", h.Wrap(h.HandleReset,
		h.LogElapsedTime,
		h.AddRequestId,
		h.FilterGet,
	))

	server := http.Server{Addr: fmt.Sprintf(":%d", h.DEPLOY_PORT)}
	go func() {
		fmt.Printf("Listening on port %d\n", h.DEPLOY_PORT)
		log.Print(server.ListenAndServe())
	}()

	quit := make(chan os.Signal)
	timeout := time.Second * 5
	signal.Notify(quit, os.Interrupt, os.Kill)
	log.Printf("got signal %s", (<-quit).String())
	log.Printf("waiting %v for graceful shutgown", timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer database.Close()
	defer controller.Close()
	defer cancel()
	log.Printf("error: %v", server.Shutdown(ctx))
}

package main

import (
	"fmt"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {
	http.Handle("/take", handler.GetHandler(http.HandlerFunc(handler.HandleTake)))
	http.Handle("/fund", handler.GetHandler(http.HandlerFunc(handler.HandleFund)))
	http.Handle("/balance", handler.GetHandler(http.HandlerFunc(handler.HandleBalance)))
	http.Handle("/announceTournament", handler.GetHandler(http.HandlerFunc(handler.HandleAnnounce)))
	http.Handle("/joinTournament", handler.GetHandler(http.HandlerFunc(handler.HandleJoin)))
	http.Handle("/resultTournament", handler.PostJsonHandler(http.HandlerFunc(handler.HandleResult)))
	http.Handle("/reset", handler.GetHandler(http.HandlerFunc(handler.HandleReset)))

	fmt.Printf("Listening on port %d\n", handler.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", handler.DEPLOY_PORT), nil)
}

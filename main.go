package main

import (
	"fmt"
	"net/http"

	h "github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {
	http.Handle("/take", h.MethodGet(http.HandlerFunc(h.HandleTake)))
	http.Handle("/fund", h.MethodGet(http.HandlerFunc(h.HandleFund)))
	http.Handle("/balance", h.MethodGet(http.HandlerFunc(h.HandleBalance)))
	http.Handle("/announceTournament", h.MethodGet(http.HandlerFunc(h.HandleAnnounce)))
	http.Handle("/joinTournament", h.MethodGet(http.HandlerFunc(h.HandleJoin)))
	http.Handle("/resultTournament", h.MethodPost(h.CntTypeJson(http.HandlerFunc(h.HandleResult))))
	http.Handle("/reset", h.MethodGet(http.HandlerFunc(h.HandleReset)))

	fmt.Printf("Listening on port %d\n", h.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", h.DEPLOY_PORT), nil)
}

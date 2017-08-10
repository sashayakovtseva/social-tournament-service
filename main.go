package main

import (
	"fmt"
	"net/http"

	h "github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {
	http.Handle("/take", h.MethodGet(h.HandleWithErrWrap(h.HandleTake)))
	http.Handle("/fund", h.MethodGet(h.HandleWithErrWrap(h.HandleFund)))
	http.Handle("/balance", h.MethodGet(h.HandleWithErrWrap(h.HandleBalance)))
	http.Handle("/announceTournament", h.MethodGet(h.HandleWithErrWrap(h.HandleAnnounce)))
	http.Handle("/joinTournament", h.MethodGet(h.HandleWithErrWrap(h.HandleJoin)))
	http.Handle("/resultTournament", h.MethodPost(h.CntTypeJson(h.HandleWithErrWrap(h.HandleResult))))
	http.Handle("/reset", h.MethodGet(h.HandleWithErrWrap(h.HandleReset)))

	fmt.Printf("Listening on port %d\n", h.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", h.DEPLOY_PORT), nil)
}

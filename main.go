package main

import (
	"fmt"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {
	http.HandleFunc("/take", handler.GetHandler(handler.HandleTake))
	http.HandleFunc("/fund", handler.GetHandler(handler.HandleFund))
	http.HandleFunc("/balance", handler.GetHandler(handler.HandleBalance))
	http.HandleFunc("/announceTournament", handler.GetHandler(handler.HandleAnnounce))
	http.HandleFunc("/joinTournament", handler.GetHandler(handler.HandleJoin))
	http.HandleFunc("/resultTournament", handler.PostHandler(handler.HandleResult))
	http.HandleFunc("/reset", handler.GetHandler(handler.HandleReset))

	fmt.Printf("Listening on port %d\n", handler.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", handler.DEPLOY_PORT), nil)
}

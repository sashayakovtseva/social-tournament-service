package main

import (
	"fmt"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {
	http.HandleFunc("/take", handler.HandleTake)
	http.HandleFunc("/fund", handler.HandleFund)
	http.HandleFunc("/balance", handler.HandleBalance)
	http.HandleFunc("/announceTournament", handler.HandleAnnounce)
	http.HandleFunc("/joinTournament", handler.HandleJoin)
	http.HandleFunc("/reset", handler.HandleReset)

	fmt.Printf("Listening on port %d\n", handler.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", handler.DEPLOY_PORT), nil)
}

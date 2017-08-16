package main

import (
	"fmt"
	"net/http"

	h "github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {
	http.Handle("/take", h.Wrap(h.HandleTake, h.MethodGet, h.CtxWithReqId))
	http.Handle("/fund", h.Wrap(h.HandleFund, h.MethodGet, h.CtxWithReqId))
	http.Handle("/balance", h.Wrap(h.HandleBalance, h.MethodGet, h.CtxWithReqId))
	http.Handle("/announceTournament", h.Wrap(h.HandleAnnounce, h.MethodGet, h.CtxWithReqId))
	http.Handle("/joinTournament", h.Wrap(h.HandleJoin, h.MethodGet, h.CtxWithReqId))
	http.Handle("/resultTournament", h.Wrap(h.HandleJoin, h.MethodPost, h.CntTypeJson, h.CtxWithReqId))
	http.Handle("/reset", h.Wrap(h.HandleReset, h.MethodGet, h.CtxWithReqId))

	fmt.Printf("Listening on port %d\n", h.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", h.DEPLOY_PORT), nil)
}

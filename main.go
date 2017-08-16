package main

import (
	"fmt"
	"net/http"

	h "github.com/sashayakovtseva/social-tournament-service/handler"
)

func main() {
	http.Handle("/take", h.Wrap(h.HandleTake,
		h.LogElapsedTime,
		h.CtxWithReqId,
		h.MethodGet,
	))
	http.Handle("/fund", h.Wrap(h.HandleFund,
		h.LogElapsedTime,
		h.CtxWithReqId,
		h.MethodGet,
	))
	http.Handle("/balance", h.Wrap(h.HandleBalance,
		h.LogElapsedTime,
		h.CtxWithReqId,
		h.MethodGet,
	))
	http.Handle("/announceTournament", h.Wrap(h.HandleAnnounce,
		h.LogElapsedTime,
		h.CtxWithReqId,
		h.MethodGet,
	))
	http.Handle("/joinTournament", h.Wrap(h.HandleJoin,
		h.LogElapsedTime,
		h.CtxWithReqId,
		h.MethodGet,
	))
	http.Handle("/resultTournament", h.Wrap(h.HandleResult,
		h.LogElapsedTime,
		h.CtxWithReqId,
		h.CntTypeJson,
		h.MethodPost,
	))
	http.Handle("/reset", h.Wrap(h.HandleReset,
		h.LogElapsedTime,
		h.CtxWithReqId,
		h.MethodGet,
	))

	fmt.Printf("Listening on port %d\n", h.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", h.DEPLOY_PORT), nil)
}

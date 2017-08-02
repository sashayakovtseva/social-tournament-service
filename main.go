package main

import (
	"fmt"
	"net/http"

	"github.com/sashayakovtseva/social-tournament-service/view"
)

func main() {
	http.HandleFunc("/fund", view.HandleFund)
	http.HandleFunc("/balance", view.HandleBalance)

	fmt.Printf("Listening on port %d\n", view.DEPLOY_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", view.DEPLOY_PORT), nil)
}

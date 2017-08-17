package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/sashayakovtseva/social-tournament-service/handler"
)

func get(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("URL: %s\nStatus: %s\nMessage: %s\n\n", url, resp.Status, b)
	return err
}

func post(url, cntType string, body io.Reader) error {
	resp, err := http.Post(url, cntType, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("URL: %s\nStatus: %s\nMessage: %s\n\n", url, resp.Status, b)
	return err
}

func TestServer(t *testing.T) {
	http.Handle("/take", h.FilterGet(h.HandleWithErrWrap(h.HandleTake)))
	http.Handle("/fund", h.FilterGet(h.HandleWithErrWrap(h.HandleFund)))
	http.Handle("/balance", h.FilterGet(h.HandleWithErrWrap(h.HandleBalance)))
	http.Handle("/announceTournament", h.FilterGet(h.HandleWithErrWrap(h.HandleAnnounce)))
	http.Handle("/joinTournament", h.FilterGet(h.HandleWithErrWrap(h.HandleJoin)))
	http.Handle("/resultTournament", h.MethodPost(h.FilterJSON(h.HandleWithErrWrap(h.HandleResult))))
	http.Handle("/reset", h.FilterGet(h.HandleWithErrWrap(h.HandleReset)))

	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	err := get(ts.URL + "/fund?playerId=P1&points=300")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/fund?playerId=P2&points=300")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/fund?playerId=P3&points=300")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/fund?playerId=P4&points=500")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/fund?playerId=P5&points=1000")
	if err != nil {
		t.Error(err)
	}

	err = get(ts.URL + "/balance?playerId=P1")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P2")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P3")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P4")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P5")
	if err != nil {
		t.Error(err)
	}

	err = get(ts.URL + "/announceTournament?tournamentId=1&deposit=1000")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/joinTournament?tournamentId=1&playerId=P5")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/joinTournament?tournamentId=1&playerId=P1&backerId=P2&backerId=P3&backerId=P4")
	if err != nil {
		t.Error(err)
	}

	err = get(ts.URL + "/balance?playerId=P1")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P2")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P3")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P4")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P5")
	if err != nil {
		t.Error(err)
	}

	err = post(ts.URL+"/resultTournament", "application/json",
		strings.NewReader("{\"tournamentId\": \"1\", \"winners\": [{\"playerId\": \"P1\", \"prize\": 2000}]}"))
	if err != nil {
		t.Error(err)
	}

	err = get(ts.URL + "/balance?playerId=P1")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P2")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P3")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P4")
	if err != nil {
		t.Error(err)
	}
	err = get(ts.URL + "/balance?playerId=P5")
	if err != nil {
		t.Error(err)
	}

	err = get(ts.URL + "/reset")
	if err != nil {
		t.Error(err)
	}
}

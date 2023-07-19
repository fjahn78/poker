package poker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := mustMakePlayerServer(t, &store)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		want := "20"

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), want)
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		want := "10"

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), want)
	})
	t.Run("return 404 for a non existent player", func(t *testing.T) {
		request := NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response
		want := http.StatusNotFound

		AssertStatus(t, got, want)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		scores:   map[string]int{},
		winCalls: nil,
	}
	server := mustMakePlayerServer(t, &store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		player := "Pepper"
		request := NewWinPostRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusAccepted)

		AssertPlayerWin(t, &store, "Pepper")
	})
}

func TestLeague(t *testing.T) {

	t.Run("it returns league as json", func(t *testing.T) {
		wantedLeague := League{
			{
				Name: "Chleo",
				Wins: 32,
			},
			{
				Name: "Chris",
				Wins: 20,
			},
			{
				Name: "Tiest",
				Wins: 14,
			},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &store)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetLeagueFromResponse(t, response.Body)

		AssertStatus(t, response, http.StatusOK)
		AssertLeague(t, got, wantedLeague)
		AssertContentType(t, response, jsonContentType)
	})
}
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{})

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
	})
	t.Run("when we get a meesage over a websocket it is a winner of a game", func(t *testing.T) {
		store := &StubPlayerStore{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		AssertPlayerWin(t, store, winner)
	})
}

func writeWSMessage(t *testing.T, ws *websocket.Conn, winner string) {
	t.Helper()
	if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func mustDialWS(t *testing.T, wsURL string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	return ws
}

func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

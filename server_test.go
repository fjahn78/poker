package poker_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	poker "github.com/fjahn78/poker"
	"github.com/gorilla/websocket"
)

var (
	dummyGame = &GameSpy{}
)

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := mustMakePlayerServer(t, &store, dummyGame)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		want := "20"

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), want)
	})
	t.Run("returns Floyd's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		want := "10"

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), want)
	})
	t.Run("return 404 for a non existent player", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response
		want := http.StatusNotFound

		poker.AssertStatus(t, got, want)
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		Scores:   map[string]int{},
		WinCalls: nil,
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		player := "Pepper"
		request := poker.NewWinPostRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusAccepted)

		poker.AssertPlayerWin(t, &store, "Pepper")
	})
}

func TestLeague(t *testing.T) {

	t.Run("it returns league as json", func(t *testing.T) {
		wantedLeague := []poker.Player{
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

		store := poker.StubPlayerStore{League: wantedLeague}
		server := mustMakePlayerServer(t, &store, dummyGame)

		request := poker.NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := poker.GetLeagueFromResponse(t, response.Body)

		poker.AssertStatus(t, response, http.StatusOK)
		poker.AssertLeague(t, got, wantedLeague)
		poker.AssertContentType(t, response, "application/json")
	})
}
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &poker.StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response, http.StatusOK)
	})
	t.Run("start a game with 3 players and declare Ruth as the winner", func(t *testing.T) {
		game := &GameSpy{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)
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

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

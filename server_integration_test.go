package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poker "github.com/fjahn78/poker"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDB := createTempFile(t, `[]`)
	defer cleanDB()
	store, _ := poker.NewFileSystemPlayerStore(database)
	server := mustMakePlayerServer(t, store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), poker.NewWinPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewWinPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), poker.NewWinPostRequest(player))

	t.Run("get score", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewGetScoreRequest(player))
		poker.AssertStatus(t, response, http.StatusOK)

		poker.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, poker.NewLeagueRequest())
		poker.AssertStatus(t, response, http.StatusOK)

		got := poker.GetLeagueFromResponse(t, response.Body)
		want := []poker.Player{
			{
				Name: "Pepper",
				Wins: 3,
			},
		}
		poker.AssertLeague(t, got, want)
	})
}

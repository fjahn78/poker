package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDB := createTempFile(t, `[]`)
	defer cleanDB()
	store, _ := NewFileSystemPlayerStore(database)
	server := mustMakePlayerServer(t, store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), NewWinPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewWinPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewWinPostRequest(player))

	t.Run("get score", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewGetScoreRequest(player))
		AssertStatus(t, response, http.StatusOK)

		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest())
		AssertStatus(t, response, http.StatusOK)

		got := GetLeagueFromResponse(t, response.Body)
		want := []Player{
			{
				Name: "Pepper",
				Wins: 3,
			},
		}
		AssertLeague(t, got, want)
	})
}

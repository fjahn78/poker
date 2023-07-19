package poker_test

import (
	"os"
	"testing"

	poker "github.com/fjahn78/poker"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league sorted from file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
      {"Name": "Chleo", "Wins": 10},
      {"Name": "Chris", "Wins": 33}
    ]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)

		poker.AssertNoError(t, err)

		got := store.GetLeague()
		want := []poker.Player{
			{
				Name: "Chris",
				Wins: 33,
			},
			{
				Name: "Chleo",
				Wins: 10,
			},
		}
		poker.AssertLeague(t, got, want)

		//* read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
      {"Name": "Chleo", "Wins": 10},
      {"Name": "Chris", "Wins": 33}
    ]`)
		defer cleanDatabase()

		store, _ := poker.NewFileSystemPlayerStore(database)

		got := store.GetPlayerScore("Chris")
		want := 33

		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
      {"Name": "Chleo", "Wins": 10},
      {"Name": "Chris", "Wins": 33}
    ]`)
		defer cleanDatabase()

		store, _ := poker.NewFileSystemPlayerStore(database)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
      {"Name": "Chleo", "Wins": 10},
      {"Name": "Chris", "Wins": 33}
    ]`)
		defer cleanDatabase()

		store, _ := poker.NewFileSystemPlayerStore(database)

		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		poker.AssertScoreEquals(t, got, want)
	})
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)
	})
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	_, err = tmpfile.Write([]byte(initialData))
	if err != nil {
		t.Fatalf("could not write to temp file: %v", err)
	}

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}
	return tmpfile, removeFile
}

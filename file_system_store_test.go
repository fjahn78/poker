package main

import (
	"io"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
      {"Name": "Chleo", "Wins": 10},
      {"Name": "Chris", "Wins": 33}
    ]`)
    defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()
		want := []Player{
			{
				Name: "Chleo",
				Wins: 10,
			},
			{
				Name: "Chris",
				Wins: 33,
			},
		}
		assertLeague(t, got, want)
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
      {"Name": "Chleo", "Wins": 10},
      {"Name": "Chris", "Wins": 33}
    ]`)
    defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEquals(t, got, want)
	})
}

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
  t.Helper()

  tmpfile, err := os.CreateTemp("","db")

  if err != nil {
    t.Fatalf("could not create temp file %v", err)
  }

  _, err = tmpfile.Write([]byte(initialData))
	if err != nil {
		t.Fatalf("could not write to temp file: %s", err)
	}

  removeFile := func() {
    tmpfile.Close()
    os.Remove(tmpfile.Name())
  }
  return tmpfile, removeFile
}

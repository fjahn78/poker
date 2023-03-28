package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	_, err := f.database.Seek(0, 0)
	if err != nil {
		log.Fatal("database corrupted")
	}

	league, _ := NewLeague(f.database)
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	var wins int

	for _, player := range f.GetLeague() {
		if player.Name == name {
			wins = player.Wins
			break
		}
	}
	return wins
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()

	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
		}
	}
	// trunk-ignore(golangci-lint/errcheck)
	f.database.Seek(0, 0)
	// trunk-ignore(golangci-lint/errcheck)
	json.NewEncoder(f.database).Encode(league)
}

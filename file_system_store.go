package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league   League
}

func (f *FileSystemPlayerStore) GetLeague() League {
	_, err := f.database.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	league, _ := NewLeague(f.database)
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	player := f.GetLeague().Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	_, err := f.database.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(f.database).Encode(f.league)
	if err != nil {
		log.Fatal(err)
	}
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	if _, err := database.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database: database,
		league:   league,
	}
}

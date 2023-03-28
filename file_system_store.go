package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
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
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	_, err := f.database.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(f.database).Encode(league)
	if err != nil {
		log.Fatal(err)
	}
}
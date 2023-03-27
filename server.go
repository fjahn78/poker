package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch player {
	case "Pepper":
		fmt.Fprint(w, "20")
	case "Floyd":
		fmt.Fprint(w, "10")
	}
}

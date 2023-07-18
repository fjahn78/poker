package poker_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	poker "github.com/fjahn78/poker"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("record Chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)
		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, dummyPlayerStore, "Chris")
	})
	t.Run("record Cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)
		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, dummyPlayerStore, "Cleo")
	})
	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		got := stdout.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []poker.ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.Alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.Alerts)
				}
				got := blindAlerter.Alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})
}


func assertScheduledAlert(t testing.TB, got, want poker.ScheduledAlert) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got amount %q, want %q", got, want)
	}
}

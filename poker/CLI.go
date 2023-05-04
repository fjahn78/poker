package poker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	blindTimeIncrement = 10 * time.Minute
	PlayerPrompt       = "Please enter the number of players: "
)

type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	out         io.Writer
	alerter     BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		out:         out,
		alerter:     alerter,
	}
}

func (cli *CLI) readline() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	cli.scheduleBlindAlerts()
	userInput := cli.readline()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindTimeIncrement
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

package poker

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

// dummies

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

// helpers

func checkSchedulingCases(t *testing.T, cases []ScheduledAlert, alerter *SpyBlindAlerter) {
	t.Helper()

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(alerter.Alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, alerter)
			}

			got := alerter.Alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}

func assertScheduledAlert(t testing.TB, got, want ScheduledAlert) {
	if got.Amount != want.Amount {
		t.Errorf("got amount %d, want %d", got.Amount, want.Amount)
	}

	if got.At != want.At {
		t.Errorf("got scheduled time of %v, want %v", got.At, want.At)
	}
}

// tests

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(5, io.Discard)

		cases := []ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Minute, Amount: 200},
			{At: 20 * time.Minute, Amount: 300},
			{At: 30 * time.Minute, Amount: 400},
			{At: 40 * time.Minute, Amount: 500},
			{At: 50 * time.Minute, Amount: 600},
			{At: 60 * time.Minute, Amount: 800},
			{At: 70 * time.Minute, Amount: 1000},
			{At: 80 * time.Minute, Amount: 2000},
			{At: 90 * time.Minute, Amount: 4000},
			{At: 100 * time.Minute, Amount: 8000},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(7, io.Discard)

		cases := []ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

}

func TestGame_Finish(t *testing.T) {
	store := &StubPlayerStore{}
	game := NewTexasHoldem(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	assertPlayerWin(t, store, winner)
}

package main

import "sync"

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{scores: map[string]int{}}
}

type InMemoryPlayerStore struct {
	scores map[string]int
	mu     sync.Mutex
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.scores[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.scores[name]++
}

func (i *InMemoryPlayerStore) GetLeague() []Player {
	var league []Player
	for name, wins := range i.scores {
		league = append(league, Player{name, wins})
	}

	return league
}

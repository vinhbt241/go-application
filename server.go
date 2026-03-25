package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

const jsonContentType = "application/json"
const htmlTemplateFile = "game.html"

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
	game     Game
}

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)

	htmlTemplatePath := filepath.Join(baseDir, htmlTemplateFile)
	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store
	p.game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.playGame))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

// helpers

func (server *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := server.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (server *PlayerServer) processWin(w http.ResponseWriter, player string) {
	server.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (server *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(server.store.GetLeague())
}

func (server *PlayerServer) getLeagueTable() League {
	return League{
		{"Chris", 20},
	}
}

func (server *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		server.processWin(w, player)
	case http.MethodGet:
		server.showScore(w, player)
	}
}

func (p *PlayerServer) playGame(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.game.Start(numberOfPlayers, io.Discard) //todo: Don't discard the blinds messages!

	_, winner, _ := conn.ReadMessage()
	p.game.Finish(string(winner))
}

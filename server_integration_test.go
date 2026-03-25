package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThemForInMemoryStore(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := mustMakePlayerServer(t, store, dummygame)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		assertStatus(t, response, http.StatusOK)

		got := getLeagueFromReponse(t, response.Body)
		want := League{{"Pepper", 3}}
		assertLeague(t, got, want)
	})
}

func TestRecordingWinsAndRetrievingThemForPostgresStore(t *testing.T) {
	store := NewPostgresPlayerStore(true)
	defer store.db.Close()

	server := mustMakePlayerServer(t, store, dummygame)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		assertStatus(t, response, http.StatusOK)

		got := getLeagueFromReponse(t, response.Body)
		want := League{{"Pepper", 3}}
		assertLeague(t, got, want)
	})
}

func TestRecordingWinsAndRetrievingThemForFileSystemStore(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "")
	defer cleanDatabase()
	store, _ := NewFileSystemPlayerStore(database)

	server := mustMakePlayerServer(t, store, dummygame)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		assertStatus(t, response, http.StatusOK)

		got := getLeagueFromReponse(t, response.Body)
		want := League{{"Pepper", 3}}
		assertLeague(t, got, want)
	})
}

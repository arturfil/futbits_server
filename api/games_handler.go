package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET/games
func (app *application) GetAllGames(w http.ResponseWriter, r *http.Request) {
	var games data.Game
	all, err := games.GetAllGames()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"games": all})
}

// GET/game/id
func (app *application) GetGameById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	game, err := app.models.Game.GetGameById(id)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, game)
}

// POST/game
func (app *application) CreateGame(w http.ResponseWriter, r *http.Request) {
	var g data.Game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusOK, g)
	id, err := app.models.Game.CreateGame(g)
	if err != nil {
		app.errorLog.Println(err)
		newGame, _ := app.models.Game.GetGameById(id)
		app.writeJSON(w, http.StatusOK, newGame)
	}
}

// PUT/game
func (app *application) UpdateGame(w http.ResponseWriter, r *http.Request) {
	var g data.Game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusOK, g)
	err = g.UpdateGame()
	if err != nil {
		app.errorLog.Println(err)
		app.writeJSON(w, http.StatusOK, "Updated Game")
	}
}

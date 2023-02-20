package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET/games
func GetAllGames(w http.ResponseWriter, r *http.Request) {
	var games models.Game
	all, err := games.GetAllGames()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"games": all})
}

// GET/game/id
func GetGameById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	game, err := mod.Game.GetGameById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, game)
}

// POST/game
func CreateGame(w http.ResponseWriter, r *http.Request) {
	var g models.Game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, g)
	id, err := mod.Game.CreateGame(g)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		newGame, _ := mod.Game.GetGameById(id)
		helpers.WriteJSON(w, http.StatusOK, newGame)
	}
}

// PUT/game
func UpdateGame(w http.ResponseWriter, r *http.Request) {
	var g models.Game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, g)
	err = g.UpdateGame()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusOK, "Updated Game")
	}
}

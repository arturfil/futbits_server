package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET/games
func GetAllGames(w http.ResponseWriter, r *http.Request) {
	var games services.Game
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

// POST/games/game/byDateField
func GetGameByDateField(w http.ResponseWriter, r *http.Request) {
    var g services.Game
    err := json.NewDecoder(r.Body).Decode(&g)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // helpers.WriteJSON(w, http.StatusOK, g)
    if err != nil {
        helpers.MessageLogs.ErrorLog.Println(err)
    }
    game, err := mod.Game.GetGameByDateField(g)
    helpers.WriteJSON(w, http.StatusOK, game)
}

// POST/game
func CreateGame(w http.ResponseWriter, r *http.Request) {
	var g services.Game
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, g)
	user, err := mod.Game.CreateGame(g)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		newGame, _ := mod.Game.GetGameById(user.ID)
		helpers.WriteJSON(w, http.StatusOK, newGame)
	}
}

// PUT/game
func UpdateGame(w http.ResponseWriter, r *http.Request) {
	var g services.Game
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

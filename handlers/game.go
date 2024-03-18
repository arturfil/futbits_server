package handlers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TODO, check if you can delete the struct declarations
// inside the functions
var game services.Game
var h helpers.Message

// GET/games
func GetAllGames(w http.ResponseWriter, r *http.Request) {

	group_id := chi.URLParam(r, "user_id")
	log.Println("user_id", group_id)

	var games services.Game

	all, err := games.GetAllGames(group_id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"games": all})
}

// GET/game/id
func GetGameById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	game, err := game.GetGameById(id)
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
	game, err := g.GetGameByDateField(g)
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
	user, err := g.CreateGame(g)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		newGame, _ := g.GetGameById(user.ID)
		helpers.WriteJSON(w, http.StatusOK, newGame)
	}
}

// PUT/game
func UpdateGame(w http.ResponseWriter, r *http.Request) {
	var g services.Game

	id := chi.URLParam(r, "id")
	log.Println("id", id)

	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, g)

	err = g.UpdateGame(id, g)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}

	helpers.WriteJSON(w, http.StatusNoContent, "Updated Game")
}

func DeleteGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := game.DeleteGame(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	// TODO: return 203 -> successfully deleted
	helpers.WriteJSON(w, http.StatusNoContent, "successfully deleted")
}

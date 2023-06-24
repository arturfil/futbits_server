package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET/groups
func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var groups services.Group
	all, err := groups.GetAllGroups()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"groups": all})
}

// GET/groups/:user_id
func GetAllGroupsOfAUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	groups, err := mod.Group.GetGroupsByMemberId(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, groups)
}

// GET/groups/group/:id
func GetGroupById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	group, err := mod.Group.GetGroupById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, group)
}

// POST/groups/create
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var g services.Group 
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := mod.Group.CreateGroup(g)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}
	newGroup, _ := mod.Group.GetGroupById(id)
	helpers.WriteJSON(w, http.StatusOK, newGroup)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    fmt.Println("%v", id)
}

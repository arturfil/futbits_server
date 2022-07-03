package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type responseObj struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

func (app *application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users data.User
	all, err := users.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"users": all})
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var myKey = []byte(os.Getenv("SECRET_KEY"))
	type credentials struct {
		UserName string `json:"email"`
		Password string `json:"password"`
	}
	var creds credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "Invalid json supplied"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}
	// get user if creds are valid
	user, err := app.models.User.GetByEmail(creds.UserName)
	if err != nil {
		app.errorJSON(w, errors.New("invalid no user found"))
		return
	}
	// check if valid
	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		app.errorJSON(w, errors.New("invalid username/password"))
		return
	}
	// create new token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = user.FirstName
	claims["exp"] = time.Now().Add(time.Minute * 60 * 4).Unix()

	tokenString, err := token.SignedString(myKey)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}
	// create response
	response := tokenResponse{
		Token: tokenString,
	}

	// send response if no erros
	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Signup(w http.ResponseWriter, r *http.Request) {
	var u data.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusOK, u)
	id, err := app.models.User.Signup(u)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err, http.StatusForbidden)
		app.infoLog.Println("Got back if of", id)
		newUser, _ := app.models.User.GetById(id)
		app.writeJSON(w, http.StatusOK, newUser)
	}
}

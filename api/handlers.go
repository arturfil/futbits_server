package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"errors"
	"fmt"
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
	Token string     `json:"token"`
	User  *data.User `json:"user"`
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
	claims["name"] = user.FirstName
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 60 * 4).Unix()

	tokenString, err := token.SignedString(myKey)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}
	user.Password = "hidden"
	// create response
	response := tokenResponse{
		Token: tokenString,
		User:  user,
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

func (app *application) SearchUser(w http.ResponseWriter, r *http.Request) {
	searchWord := r.URL.Query().Get("keyword")
	app.infoLog.Println(searchWord)
	users, err := app.models.User.SearchUserBy(searchWord)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"users": users})

}

func (app *application) GetUserByToken(w http.ResponseWriter, r *http.Request) {
	type TokenClaim struct {
		Authorized bool   `json:"authorized"`
		Email      string `json:"email"`
		Exp        int    `json:"exp"`
		Name       string `json:"name"`
		jwt.StandardClaims
	}

	claims := &TokenClaim{}
	token, err := jwt.ParseWithClaims(r.Header["Authorization"][0], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return myKey, nil
	})
	if err != nil {
		app.infoLog.Print(err)
		return
	}
	if token.Valid {

		if err != nil {
			app.infoLog.Print(err)
		}

		if err != nil {
			app.infoLog.Print(err)
			return
		}
		user, err := app.models.User.GetByEmail(claims.Email)
		if err != nil {
			app.infoLog.Print(err)
		}
		user.Password = ""
		app.writeJSON(w, http.StatusOK, user)
		return
	}
}

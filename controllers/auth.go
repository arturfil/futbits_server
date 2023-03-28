package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var u services.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, u)
	_, err = mod.User.Signup(u)
	if err != nil {
		h.ErrorLog.Panicln(err)
		helpers.ErrorJSON(w, err, http.StatusForbidden)
	}
}

// this method will check that the credentials against the key are
// equal
func Login(w http.ResponseWriter, r *http.Request) {
	var myKey = []byte(os.Getenv("SECRET_KEY"))
	type credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// setup creds & jsonResponse
	var creds credentials
	var payload services.JsonResponse

	// read Json
	err := helpers.ReadJson(w, r, &creds)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "Invalid json supplied"
		_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
	}
	// get user if creds are valid
	user, err := mod.User.GetByEmail(creds.Email)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid, no user found"))
		return
	}
	// check if valid
	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		helpers.ErrorJSON(w, errors.New("invalid email/password"))
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
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err)
		return
	}
	user.Password = "hidden"
	// create response
	response := services.TokenResponse{
		Token: tokenString,
		User:  user,
	}
	// send response if no errors
	err = helpers.WriteJSON(w, http.StatusOK, response)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}
}

// GetUserByToken - this method will serve to get the user once a token is provided. This
// is really usefull when ever you want to preserve login state in the front-end
func GetUserByToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached [GetUserByToken] method")
	var myKey = []byte(os.Getenv("SECRET_KEY"))
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
	if token.Valid {
		if err != nil {
			helpers.MessageLogs.InfoLog.Print(err)
		}

		if err != nil {
			helpers.MessageLogs.InfoLog.Print(err)
			return
		}

		user, err := mod.User.GetByEmail(claims.Email)
		if err != nil {
			helpers.MessageLogs.InfoLog.Print(err)
		}
		user.Password = ""
		helpers.WriteJSON(w, http.StatusOK, user)
		return
	}
}

// search user by keyword | name | lastName | email
func SearchUser(w http.ResponseWriter, r *http.Request) {
	searchWord := r.URL.Query().Get("keyword")
	helpers.MessageLogs.InfoLog.Println(searchWord)
	users, err := mod.User.SearchUserBy(searchWord)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"users": users})
}

// GET user by id
func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := mod.User.GetUserById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, user)
}

// this function resets the user password
// func (u *User) ResetPassword(password string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
// 	if err != nil {
// 		return err
// 	}

// 	stmt := `update users set password = $1 where id = $2`
// 	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

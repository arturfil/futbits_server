package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

// time for db process any transaction
const dbTimeout = time.Second * 3

var db *sql.DB

// create a new db pool to cache data
func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User:  User{},
		Token: Token{},
	}
}

// models we will use, User & Token
type Models struct {
	User   User
	Field  Field
	Token  Token
	Game   Game
	Group  Group
	Member Member
}

// User model
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// this function will return all the users
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, created_at, updated_at from users order by last_name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// this function will return a user, given the email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where email = $1`
	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// this function will return a user, given the id
func (u *User) GetById(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where id = $1`
	var user User

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// this will help when creating a new member for a group
func (u *User) SearchUserBy(term string) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select 
			id, 
			email, 
			first_name, 
			last_name, 
			created_at, 
			updated_at 
		from users 
		where first_name like '%' || $1 || '%' or last_name like '%' || $1 || '%' or email like '%' || $1 || '%'
	`
	rows, err := db.QueryContext(ctx, query, term)
	if err != nil {
		fmt.Println("HERE", err)
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// this function will update a user
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		updated_at = $4
		where id = $5
	`

	_, err := db.ExecContext(
		ctx,
		stmt,
		u.Email,
		u.LastName,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// this function will delete a user
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, stmt, u.ID)

	if err != nil {
		return err
	}
	return nil
}

// this function will create a user
func (u *User) Signup(user User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return "0", err
	}

	newId := uuid.New()
	stmt := `
		insert into users (id, email, first_name, last_name, password, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id
	`
	err = db.QueryRowContext(ctx, stmt,
		newId,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return "0", err
	}
	return newId.String(), nil
}

// this function resets the user password
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// This function will check whether the passwords match or not
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// token model
type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Expiry    time.Time `json:"expiry"`
}

// Gets the user's TOKEN by plain text token
func (t *Token) GetTokenByToken(plainText string) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, user_id, email, token, token_hash, created_at, updated_at, expiry
		from tokens where token = $1
	`

	var token Token
	row := db.QueryRowContext(ctx, query, plainText)
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Email,
		&token.Token,
		&token.TokenHash,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Expiry,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Gets the USER, by token
func (t *Token) GetUserByToken(token Token) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at, from users where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, token.UserID)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Generates a token
func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]
	return token, nil
}

// make sure we have a valid token
func (t *Token) AuthenticateToken(r *http.Request) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no valid authorization headers received")
	}
	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("token is not valid")
	}
	tkn, err := t.GetTokenByToken(token)
	if err != nil {
		return nil, errors.New("no matching token found")
	}
	if tkn.Expiry.Before(time.Now()) {
		return nil, errors.New("expired token")
	}
	user, err := t.GetUserByToken(*tkn)
	if err != nil {
		return nil, errors.New("no matching user foudn")
	}
	return user, nil
}

// create token
func (t *Token) InsertToken(token Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from tokens where user_id = $1`
	_, err := db.ExecContext(ctx, stmt, token.UserID)
	if err != nil {
		return err
	}
	token.Email = u.Email
	stmt = `insert into tokens (user_id, email, token, token_hash, created_at, updated_at, expiry)
		values ($1, $2, $3, $4, $5, $6, $7)	
	`
	_, err = db.ExecContext(ctx, stmt,
		&token.UserID,
		&token.Email,
		&token.Token,
		&token.TokenHash,
		token.Expiry,
	)
	if err != nil {
		return err
	}
	return nil
}

// delete token
func (t *Token) DeleteByToken(plainText string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from tokens where token = $1`
	_, err := db.ExecContext(ctx, stmt, plainText)
	if err != nil {
		return err
	}
	return nil
}

// Validating a token
func (t *Token) ValidToken(plainText string) (bool, error) {
	token, err := t.GetTokenByToken(plainText)
	if err != nil {
		return false, errors.New("no matching token found")
	}
	_, err = t.GetUserByToken(*token)
	if err != nil {
		return false, errors.New("no matching user found")
	}
	if token.Expiry.Before(time.Now()) {
		return false, errors.New("expired token")
	}
	return true, nil
}
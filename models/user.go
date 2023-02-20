package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// method returns a user provided a valid email
// remember that here you DO need to sepcify the fields otherwise if you use *
// it wont return the result because clearance right now is NULL
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select 
			id, email, first_name, last_name, password, created_at, updated_at 
		from users 
		where email = $1
	`
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

// this query returns the user given that a vliad id is provided
func (u *User) GetUserById(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select * from users where id = $1`
	var user User
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
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

// this function will create a user
func Signup(u User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	var user User
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return nil, err
	}

	stmt := `
		insert into users (email, first_name, last_name, password, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning *;
	`
	res, err := db.ExecContext(
		ctx,
		stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		hashedPassword,
		time.Now(),
		time.Now(),
	)

	fmt.Println("RES exec -> ", res)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

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
func Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}
	return nil
}

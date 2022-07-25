package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		User: User{},
	}
}

// models we will use, User & Token
type Models struct {
	User    User
	Field   Field
	Game    Game
	Profile Profile
	Group   Group
	Member  Member
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

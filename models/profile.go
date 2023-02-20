package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Nationality string    `json:"nationality"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Position    string    `json:"position"`
	Level       string    `json:"level"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Profile) CreateProfile(profile Profile) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	newId := uuid.New()
	query := `
		insert into profile (id, user_id, nationality, age, gender, position, level, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`
	err := db.QueryRowContext(
		ctx,
		query,
		newId,
		profile.UserID,
		profile.Nationality,
		profile.Age,
		profile.Gender,
		profile.Position,
		profile.Level,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return "0", err
	}
	return newId.String(), nil
}

func (p *Profile) GetProfileById(id string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select * from profile where user_id = $1` // make sure you declare everything in scan
	var profile Profile
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Nationality,
		&profile.Age,
		&profile.Gender,
		&profile.Position,
		&profile.Level,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

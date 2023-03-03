package models

import (
	"context"
	"time"
)

type Game struct {
	ID         string    `json:"id"`
	FieldID    string    `json:"field_id"`
	StartTime  time.Time `json:"start_time"`
	MaxPlayers int8      `json:"max_players"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GameResponse struct {
	ID         string    `json:"id"`
	FieldID    string    `json:"field_id"`
	FieldName  string    `json:"field_name"`
	MaxPlayers int8      `json:"max_players"`
	StartTime  time.Time `json:"start_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (g *Game) GetAllGames() ([]*GameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// query := `
	// 	select
	// 		g.id,
	// 		g.field_id,
	// 		f.name,
	// 		g.start_time,
	// 		g.max_players,
	// 		g.created_at,
	// 		g.created_at
	// 	from
	// 		games g
	// 	inner join fields f
	// 		on g.field_id = f.id

	// `

	query := `
		select id, field_id, max_players, start_time, created_at, updated_at from games;
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var games []*GameResponse
	for rows.Next() {
		var game GameResponse
		err := rows.Scan(
			&game.ID,
			&game.FieldID,
			&game.MaxPlayers,
			&game.StartTime,
			&game.CreatedAt,
			&game.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, &game)
	}
	return games, nil
}

// GET/games/game/:id
func (g *Game) GetGameById(id string) (*Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select * from games where id = $1`
	var game Game

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&game.ID,
		&game.FieldID,
		&game.MaxPlayers,
		&game.StartTime,
		&game.CreatedAt,
		&game.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

// POST/games/create
func (g *Game) CreateGame(game Game) (*Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		insert into games (field_id, start_time, max_players, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning *
	`

	_, err := db.ExecContext(
		ctx,
		query,
		game.FieldID,
		game.StartTime,
		game.MaxPlayers,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (g *Game) UpdateGame() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		update games set
		field_id = $1,
		start_time = $2,
		max_players = $3
		updated_at = $4,
		where id = $5
	`

	_, err := db.ExecContext(
		ctx,
		query,
		g.FieldID,
		g.StartTime,
		time.Now(),
		g.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) DeleteGame() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `delete from fields where id = $1`
	_, err := db.ExecContext(ctx, query, g.ID)
	if err != nil {
		return err
	}
	return nil
}

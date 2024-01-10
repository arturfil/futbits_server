package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Game struct {
	ID        string    `json:"id"`
	FieldID   string    `json:"field_id"`
	GameDate  time.Time `json:"game_date"`
	Score     string    `json:"score"`
	GroupID   string    `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GameResponse struct {
	ID        string    `json:"id"`
	FieldID   string    `json:"field_id"`
	FieldName string    `json:"field_name"`
	GameDate  time.Time `json:"game_date"`
	GroupID   string    `json:"group_id"`
	GroupName string    `json:"group_name,omitempty"`
	Score     string    `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (g *Game) GetAllGames(user_id string) ([]*GameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
            g.id,
            g.field_id,
            f.name,
            g.score,
            g.group_id,
            g.game_date,
            g.created_at,
            g.updated_at 
        from games g 
        INNER JOIN "groups" gp ON g.group_id = gp.id 
        INNER JOIN "members" m ON gp.id = m.group_id
        INNER JOIN users u ON u.id = m.user_id 
        INNER JOIN fields f ON g.field_id = f.id
        WHERE u.id = $1;
    `
	rows, err := db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	var games []*GameResponse

	for rows.Next() {
		var game GameResponse
		err := rows.Scan(
			&game.ID,
			&game.FieldID,
			&game.FieldName,
			&game.Score,
			&game.GroupID,
			&game.GameDate,
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
func (g *Game) GetGameById(id string) (*GameResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        select 
            g.id,
            g.field_id,
            f.name,
            gp.name,
            g.group_id,
            g.score,
            g.game_date,
            g.created_at,
            g.updated_at 
        from games g
        inner join fields f on g.field_id = f.id
        inner join "groups" gp on gp.id = g.group_id
        where g.id = $1`
	var game GameResponse

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&game.ID,
		&game.FieldID,
		&game.FieldName,
		&game.GroupName,
        &game.GroupID,
		&game.Score,
		&game.GameDate,
		&game.CreatedAt,
		&game.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

// POST/games/game/byDateField
func (g *Game) GetGameByDateField(game Game) (*GameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
        select 
           g.id,                                           
            g.field_id,                                     
            f.name,                                         
            g.game_date,                                    
            g.score,                                  
            g.created_at,                                   
            g.updated_at                                    
        from games g                  
        inner join fields f                      
            on g.field_id = f.id
        where g.game_date = $1 and g.field_id = $2;    
    `
	var gameRes GameResponse

	row := db.QueryRowContext(ctx, query, game.GameDate, game.FieldID)
	err := row.Scan(
		&gameRes.ID,
		&gameRes.FieldID,
		&gameRes.FieldName,
		&gameRes.GameDate,
		&gameRes.Score,
		&gameRes.CreatedAt,
		&gameRes.UpdatedAt,
	)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, err
	}
	return &gameRes, nil
}

// POST/games/create
func (g *Game) CreateGame(game Game) (*Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	gameExists, err := g.GetGameByDateField(game)
	if err != nil {
		fmt.Println("A game with that date and time already exists") 
	}

	existsError := errors.New("Game already exists")
	fmt.Printf("Exists -> %v, aboutToCreate -> %v", gameExists, game)
	if gameExists != nil && gameExists.GameDate.UTC() == game.GameDate {
		fmt.Println("\nGame with same datetime & field already exists----->")
		return nil, existsError
	}

	query := `
		insert into games (field_id, group_id, game_date, score, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning *
	`

	_, err = db.ExecContext(
		ctx,
		query,
		game.FieldID,
		game.GroupID,
		game.GameDate,
		game.Score,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (g *Game) UpdateGame(id string, game Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    log.Println("game", game)
    
	query := `
    UPDATE games
    SET 
        field_id = $1,
        score = $2,
        game_date = $3,
        updated_at = $4
    WHERE id = $5;
			`

	_, err := db.ExecContext(
		ctx,
		query,
        game.FieldID,
        game.Score,
        game.GameDate,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) DeleteGame(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from games where id = $1`

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

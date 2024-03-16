package services

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"time"
)

// TODO: missing mapping out table to struct
type ReportDTO struct {
	ID            string         `json:"id"`
	TeamSide      string         `json:"team_side"`
	UserID        sql.NullString `json:"user_id,omitempty"`
	GameID        string         `json:"game_id"`
	PlayerName    string         `json:"player_name"`
	Goals         int            `json:"goals"`
	Assists       int            `json:"assists"`
	Won           bool           `json:"won"`
	ManOfTheMatch bool           `json:"man_of_the_match"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Report struct {
	ID            string    `json:"id"`
	TeamSide      string    `json:"team_side"`
	UserID        string    `json:"user_id,omitempty"`
	GameID        string    `json:"game_id"`
	PlayerName    string    `json:"player_name"`
	Goals         int       `json:"goals"`
	Assists       int       `json:"assists"`
	Won           bool      `json:"won"`
	ManOfTheMatch bool      `json:"man_of_the_match"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (r *Report) GetAllReports() ([]*ReportDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        select  
            id,
            team_side,
            user_id,
            game_id,
            player_name,
            goals,
            assits,
            won,
            man_of_the_match,
            created_at,
            updated_at,
        from reports
    `

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	var reports []*ReportDTO
	for rows.Next() {
		var report ReportDTO
		err := rows.Scan(
			r.ID,
			r.TeamSide,
			r.UserID,
			r.GameID,
			r.PlayerName,
			r.Goals,
			r.Assists,
			r.Won,
			r.ManOfTheMatch,
			r.CreatedAt,
			r.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}

func (r *Report) CreateReport(report Report) (*Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		insert into reports (
            team_side,
            user_id,
            game_id,
            player_name,
            goals,
            assists,
            won,
            man_of_the_match,
            created_at,
            updated_at
        )
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning *
	`
	_, err := db.ExecContext(
		ctx,
		query,
		report.TeamSide,
		report.UserID,
		report.GameID,
		report.PlayerName,
		report.Goals,
		report.Assists,
		report.Won,
		report.ManOfTheMatch,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *Report) UploadReport(file *csv.Reader) {
	fmt.Println(file)

	lines, err := file.ReadAll()
	if err != nil {
		return
	}

	for _, row := range lines {
		var query string
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

		defer cancel()

		fmt.Println(row)

		if row[2] != "" {
			query = `
                insert into reports (
                    team_side,
                    game_id,
                    user_id,
                    player_name,
                    goals,
                    assists,
                    won,
                    man_of_the_match,
                    created_at,
                    updated_at
            )
                values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning *
            `
			_, err = db.ExecContext(
				ctx,
				query,
				row[0],
				row[1],
				row[2],
				row[3],
				row[4],
				row[5],
				row[6],
				row[7],
				time.Now(),
				time.Now(),
			)
		} else {
			query = `
                insert into reports (
                    team_side,
                    game_id,
                    player_name,
                    goals,
                    assists,
                    won,
                    man_of_the_match,
                    created_at,
                    updated_at
            )
                values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *
            `
			_, err = db.ExecContext(
				ctx,
				query,
				row[0],
				row[1],
				row[3],
				row[4],
				row[5],
				row[6],
				row[7],
				time.Now(),
				time.Now(),
			)
		}

		if err != nil {
			fmt.Println("Error in insert: ", err)
		}
	}
	return
}

func (r *Report) GetReportById(id string) (*Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select * from reports where id = $1`
	var report Report

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&report.ID,
		&report.TeamSide,
		&report.UserID,
		&report.GameID,
		&report.PlayerName,
		&report.Goals,
		&report.Assists,
		&report.Won,
		&report.ManOfTheMatch,
		&report.CreatedAt,
		&report.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *Report) GetAllReportsByGroupId(group_id string) ([]*ReportDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
        select 
            r.id,
            r.team_side,
            r.user_id,
            r.game_id,
            r.player_name,
            r.goals,
            r.assists,
            r.won,
            r.man_of_the_match,
            r.created_at,
            r.updated_at
        from reports r
        inner join games g on r.game_id = g.id 
        where g.group_id = $1
        limit 40 
        offset 0;
    `
	rows, err := db.QueryContext(ctx, query, group_id)
	if err != nil {
		return nil, err
	}
	var reports []*ReportDTO
	for rows.Next() {
		var report ReportDTO

		err := rows.Scan(
			&report.ID,
			&report.TeamSide,
			&report.UserID,
			&report.GameID,
			&report.PlayerName,
			&report.Goals,
			&report.Assists,
			&report.Won,
			&report.ManOfTheMatch,
			&report.CreatedAt,
			&report.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}

func (r *Report) GetAllReportsByGameId(game_id string) ([]*ReportDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from reports where game_id = $1`

	rows, err := db.QueryContext(ctx, query, game_id)
	if err != nil {
		return nil, err
	}

	var reports []*ReportDTO
	for rows.Next() {
		var report ReportDTO
		err := rows.Scan(
			&report.ID,
			&report.TeamSide,
			&report.UserID,
			&report.GameID,
			&report.PlayerName,
			&report.Goals,
			&report.Assists,
			&report.Won,
			&report.ManOfTheMatch,
			&report.CreatedAt,
			&report.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		reports = append(reports, &report)
	}
	return reports, nil
}

func (r *Report) GetAllReporstByUserId(user_id string) ([]*Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select * from reports where user_id = $1
	`
	rows, err := db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	var reports []*Report
	for rows.Next() {
		var report Report
		err := rows.Scan(
			&report.ID,
			&report.TeamSide,
			&report.UserID,
			&report.GameID,
			&report.PlayerName,
			&report.Goals,
			&report.Assists,
			&report.Won,
			&report.ManOfTheMatch,
			&report.CreatedAt,
			&report.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}

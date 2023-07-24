package services

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	GameID        string    `json:"game_id"`
	Assists       int       `json:"assists"`
	Goals         int       `json:"goals"`
	Attendance    int       `json:"attendance"`
	ManOfTheMatch int       `json:"man_of_the_match"`
	Involvement   int       `json:"involvement"`
	Attitude      string    `json:"attitude"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (r *Report) GetAllReports() ([]*Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from reports`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var reports []*Report
	for rows.Next() {
		var report Report
		err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.GameID,
			&report.Assists,
			&report.ManOfTheMatch,
			&report.Attendance,
			&report.Involvement,
			&report.Attitude,
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

func (r *Report) CreateReport(report Report) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	newId := uuid.New()
	query := `
		insert into reports (id, user_id, game_id, assists, goals, attendance, man_of_the_match, involvement, attitude, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id
	`
	err := db.QueryRowContext(
		ctx,
		query,
		newId,
		report.UserID,
		report.GameID,
		report.Assists,
		report.Goals,
		report.Attendance,
		report.ManOfTheMatch,
		report.Involvement,
		report.Attitude,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return "0", err
	}
	return newId.String(), nil
}

func (r *Report) GetReportById(id string) (*Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select * from reports where id = $1`
	var report Report

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&report.ID,
		&report.UserID,
		&report.GameID,
		&report.Assists,
		&report.Goals,
		&report.Attendance,
		&report.ManOfTheMatch,
		&report.Involvement,
		&report.Attitude,
		&report.CreatedAt,
		&report.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *Report) GetAllReporstById(user_id string) ([]*Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select
			r.id,
			r.user_id,
			r.game_id,
			r.assists,
			r.goals,
			r.attendance,
			r.man_of_the_match,
			r.involvement,
			r.attitude,
			r.created_at,
			r.updated_at
		from
			reports r
		inner join users u on u.id = r.user_id
		where u.id = $1
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
			&report.UserID,
			&report.GameID,
			&report.Assists,
			&report.Goals,
			&report.Attendance,
			&report.ManOfTheMatch,
			&report.Involvement,
			&report.Attitude,
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

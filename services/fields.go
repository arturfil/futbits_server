package services

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FieldRepo interface {
	GetAllFields() ([]Field, error)
	GetFieldById(id string) (*Field, error)
	CreateField(field Field) (string, error)
	UpdateField() error
	DeleteField() error
}

type Field struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GET/allFields
func (f *Field) GetAllFields() ([]Field, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from fields`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var fields []Field
	for rows.Next() {
		var field Field
		err := rows.Scan(
			&field.ID,
			&field.Name,
			&field.Address,
			&field.CreatedAt,
			&field.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

// GET/fields/field/:id
func (f *Field) GetFieldById(id string) (*Field, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select * from fields where id = $1`
	var field Field

	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&field.ID,
		&field.Name,
		&field.Address,
		&field.CreatedAt,
		&field.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &field, nil
}

// POST/createField
func (f *Field) CreateField(field Field) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	newId := uuid.New()
	query := `
		insert into fields (id, name, address, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning id
	`

	err := db.QueryRowContext(ctx, query,
		newId,
		field.Name,
		field.Address,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return "0", err
	}
	return newId.String(), nil
}

// PUT/games/game
func (f *Field) UpdateField() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		update fields set
		name = $1,
		address = $2,
		updated_at = $3
		where id = $4
	`

	_, err := db.ExecContext(
		ctx,
		query,
		f.Name,
		f.Address,
		time.Now(),
		f.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (f *Field) DeleteField() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `delete from fields where id = $1`
	_, err := db.ExecContext(ctx, query, f.ID)
	if err != nil {
		return err
	}
	return nil
}

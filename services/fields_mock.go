package services

import "time"

type FieldsMock struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *FieldsMock) GetAllFields() ([]FieldsMock, error) {
    fields := []FieldsMock{
        {
            ID: "47a4a57f-fa22-4444-8f56-2a6d92f557cc",
            Name:"Test1",
            Address :"Addres Test",
            CreatedAt: time.Date(2024, 2, 16, 24, 0, 0, 0, time.UTC),
            UpdatedAt: time.Date(2024, 2, 16, 24, 0, 0, 0, time.UTC),
        },
    }
	return fields, nil
}

// GET/fields/field/:id
func (f *FieldsMock) GetFieldById(id string) (*Field, error) {
    panic("Missing implementint")
}

// POST/createFieldsMock
func (f *FieldsMock) CreateField(field Field) (error) {
    panic("Missing implementint")
}

// PUT/games/game
func (f *FieldsMock) UpdateField() error {
    panic("Missing implementint")
}

func (f *FieldsMock) DeleteField() error {
    panic("Missing implementint")
}

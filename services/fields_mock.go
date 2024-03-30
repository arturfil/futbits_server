package services

import (
	"time"
)

type MockField struct {
}

func (f *MockField) GetAllFields() ([]Field, error) {
	// fmt.Println("MockField reached")
	// fields := []Field{
	//     {
	//         ID: "47a4a57f-fa22-4444-8f56-2a6d92f557cc",
	//         Name:"Test1",
	//         Address :"Addres Test",
	//         CreatedAt: time.Date(2024, 2, 16, 24, 0, 0, 0, time.UTC),
	//         UpdatedAt: time.Date(2024, 2, 16, 24, 0, 0, 0, time.UTC),
	//     },
	// }
	return nil, nil
}

// GET/fields/field/:id
func (f *MockField) GetFieldById(id string) (*Field, error) {
	field := Field{
		ID:        "47a4a57f-fa22-4444-8f56-2a6d92f557cc",
		Name:      "Test1",
		Address:   "Addres Test",
		CreatedAt: time.Date(2024, 2, 16, 24, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 2, 16, 24, 0, 0, 0, time.UTC),
	}
	return &field, nil

}

// POST/createMockField
func (f *MockField) CreateField(field Field) error {
	panic("Missing implementint")
}

// PUT/games/game
func (f *MockField) UpdateField() error {
	panic("Missing implementint")
}

func (f *MockField) DeleteField() error {
	panic("Missing implementint")
}

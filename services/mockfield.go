package services

import "time"

type MockField struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateField implements FieldRepo.
func (f *MockField) CreateField(field Field) (string, error) {
	panic("unimplemented")
}

// DeleteField implements FieldRepo.
func (f *MockField) DeleteField() error {
	panic("unimplemented")
}

// GetAllFields implements FieldRepo.
func (f *MockField) GetAllFields() ([]Field, error) {
	var fields []Field
	return fields, nil
}

// GetFieldById implements FieldRepo.
func (f *MockField) GetFieldById(id string) (*Field, error) {
	panic("unimplemented")
}

// UpdateField implements FieldRepo.
func (f *MockField) UpdateField() error {
	panic("unimplemented")
}

// GET/allFields
// func (f *MockField) GetAllFields() ([]Field, error) {
// }

// // GET/fields/field/:id
// func (f *Field) GetFieldById(id string) (*Field, error) {
// }
//
// // POST/createField
// func (f *Field) CreateField(field Field) (string, error) {
// }
//
// // PUT/games/game
// func (f *Field) UpdateField() error {
// }
//
// func (f *Field) DeleteField() error {
// }

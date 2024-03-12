package services

type FieldRepo interface {
	GetAllFields() ([]Field, error)
	GetFieldById(id string) (*Field, error)
	CreateField(field Field) (string, error)
	UpdateField() error
	DeleteField() error
}

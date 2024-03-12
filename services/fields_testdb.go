package services

type TestField struct {

}

func (f *TestField) GetAllFields() ([]Field, error) {
	var fields []Field
	return fields, nil
}

func (f *TestField) GetFieldById(id string) (*Field, error) {
	return nil, nil
}

// POST/createField
func (f *TestField) CreateField(field Field) (string, error) {
    return "", nil
}

// PUT/games/game
func (f *TestField) UpdateField() error {
    return nil
}

func (f *TestField) DeleteField() error {
    return nil
}

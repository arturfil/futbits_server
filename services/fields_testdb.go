package services

type TestField struct {

}

func (f *TestField) GetAllFields() ([]Field, error) {
	var fields []Field
	return fields, nil
}

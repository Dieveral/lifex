package Entities

type EntityManager interface {
	Exists(key interface{}) bool
	Add(data interface{}) (int64, error)
	ShowById(id int64) error
	Show(filter interface{}) error
	ShowAll() error
	PrintInfo(entity interface{}, columnWidth []int)
}

type DefaultManager struct {
}

func (d *DefaultManager) Exists(key interface{}) bool {
	return false
}

func (d *DefaultManager) Add(data interface{}) (int64, error) {
	return 0, nil
}

func (d *DefaultManager) ShowById(id int64) error {
	return nil
}

func (d *DefaultManager) Show(filter interface{}) error {
	return nil
}

func (d *DefaultManager) ShowAll() error {
	return nil
}

func (d *DefaultManager) PrintInfo(entity interface{}, columnWidth []int) {

}

func (d *DefaultManager) PrintHeader(columnName []string, columnWidth []int) {

}

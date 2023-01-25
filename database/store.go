package database

type Store struct {
}

func New() (*Store, error) {
	store := &Store{}
	return store, nil
}

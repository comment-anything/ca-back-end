package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/comment-anything/ca-back-end/config"
	"github.com/comment-anything/ca-back-end/database/generated"
	_ "github.com/lib/pq"
)

// Store wraps the queries object, connecting it to the configured database.
type Store struct {
	DB      *sql.DB
	Queries *generated.Queries
}

// New instantiates and returns a new Store. If connect is true, it also connects that store to the database using the parameters from env vars in package config.go
func New(connect bool) (*Store, error) {
	store := &Store{}
	store.DB = nil
	store.Queries = nil
	if connect {
		err := store.Connect()
		if err != nil {
			return nil, err
		} else {
			return store, err
		}
	}
	return store, nil
}

// Connect uses environment variables configured in a secret .env file and stored in the global config singleton to connect to the Postgres server the given port.
func (s *Store) Connect() error {
	if !config.Vals.IsLoaded {
		return errors.New("Can't connect to database: environment variables aren't loaded.")
	}
	if s.DB != nil {
		s.Disconnect()
	}
	postgres, err := sql.Open("postgres", config.Vals.DB.ConnectString())

	if err != nil {
		return err
	}
	fmt.Printf("\tConnecting to database: %s\n", config.Vals.DB.DBname)
	if err != nil {
		return err
	} else {
		s.DB = postgres
		s.Queries = generated.New(s.DB)
		return nil
	}
}

// Disconnect disconnects Store from the Postgres instance.
func (s *Store) Disconnect() error {
	err := s.DB.Close()
	s.DB = nil
	s.Queries = nil
	return err
}

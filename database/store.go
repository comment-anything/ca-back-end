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
			fmt.Printf("\nError connecting to database! %s", err.Error())
			return nil, err
		} else {
			return store, nil
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

	var cstring string
	if config.Vals.DockerMode {
		cstring = config.Vals.DB.ConnectString2()
	} else {
		cstring = config.Vals.DB.ConnectString1()

	}
	postgres, err := sql.Open("postgres", cstring)
	if err != nil {
		fmt.Printf("\nFailed to open postgres! %s", err.Error())
	}

	err = postgres.Ping()
	if err != nil {
		var s string
		if config.Vals.DockerMode {
			s = config.Vals.DB.ConnectString2()
		} else {
			s = config.Vals.DB.ConnectString1()
		}
		fmt.Printf("\nFailed to ping DB with %s! %s", s, err.Error())
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

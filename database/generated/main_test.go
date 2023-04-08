package generated

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/comment-anything/ca-back-end/config"

	_ "github.com/lib/pq"
)

/*
Package generated consists of sqlc generated code for interacting with the database.
This file, main_test.go, is run when conducting unit tests of the generated code.
It creates a database connection for use by other tests.
*/

// testQueries is the instance of the query struct used by generated unit tests.
var testQueries *Queries

// TestMain is ran before all other tests in the generated package. It initializes the testDb and testQueries objects for use by other tests. It uses the values in util.Config to connect to the database; if the environment variables are not configured properly, the tests cannot run.
func TestMain(m *testing.M) {

	config.Vals.Load("../../.env", false)

	postgres, err := sql.Open("postgres", config.Vals.DB.ConnectString1())
	if err != nil {
		fmt.Println(" Couldn't run tests: Database failed to load.")
		log.Fatal(err)
	} else {
		testQueries = New(postgres)
		os.Exit(m.Run())
	}
}

package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var testEnvFile = "test.env"

func TestLoadDbEnv(t *testing.T) {
	createEnvFile()
	godotenv.Load(testEnvFile)
	Vals.loadDBEnv()
	if Vals.DB.DBname != "comm-anything-tests" {
		t.Errorf("DB_Name not correct: %v", Vals.DB.DBname)
	}
	if Vals.DB.Host != "localhost" {
		t.Errorf("DB_HOST not correct: %v", Vals.DB.DBname)
	}
	if Vals.DB.Password != "dbsuperuser991" {
		t.Errorf("DB_PASSWORD not correct: %v", Vals.DB.DBname)
	}
	if Vals.DB.Port != "5433" {
		t.Errorf("DB_HOST_PORT not correct: %v", Vals.DB.DBname)
	}
	if Vals.DB.User != "root" {
		t.Errorf("DB_USER not correct: %v", Vals.DB.DBname)
	}
	deleteEnvFile()
}

func TestLoad(t *testing.T) {
	createEnvFile()
	Vals.Load(testEnvFile)
	deleteEnvFile()
}

func TestLoadServerEnv(t *testing.T) {
	createEnvFile()
	godotenv.Load(testEnvFile)
	Vals.loadServerEnv()
	if Vals.Server.DoesLogAll != true {
		t.Errorf("SERVER_LOG_ALL is not correct: %v", Vals.Server.DoesLogAll)
	}
	if Vals.Server.JWTCookieName != "canywauth" {
		t.Errorf("JWT_COOKIE_NAME is not correct: %v", Vals.Server.JWTCookieName)
	}
	if Vals.Server.JWTKey != "key1111111" {
		t.Errorf("JWT_COOKIE_NAME is not correct: %v", Vals.Server.JWTKey)
	}
	if Vals.Server.Port != ":3000" {
		t.Errorf("SERVER_PORT is not correct: %v", Vals.Server.Port)
	}
	deleteEnvFile()
}

func createEnvFile() {
	f, e := os.Create(testEnvFile)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	fmt.Fprintln(f, "CA_TESTING_MODE=true")
	fmt.Fprintln(f, "SERVER_LOG_ALL=true")
	fmt.Fprintln(f, "DB_IMAGE=postgres:14.5-alpine")
	fmt.Fprintln(f, "DB_CONTAINER_NAME=923postgres")
	fmt.Fprintln(f, "DB_CONTAINER_PORT=5432")
	fmt.Fprintln(f, "DB_HOST=localhost")
	fmt.Fprintln(f, "DB_HOST_PORT=5433")
	fmt.Fprintln(f, "DB_USER=root")
	fmt.Fprintln(f, "DB_PASSWORD=dbsuperuser991")
	fmt.Fprintln(f, "SERVER_PORT = 3000  ")
	fmt.Fprintln(f, "JWT_KEY=key1111111")
	fmt.Fprintln(f, "JWT_COOKIE_NAME= canywauth")
	fmt.Fprintln(f, "DB_DATABASE_NAME=comm-anything")
	fmt.Fprintln(f, "TEST_DB_DATABASE_NAME = comm-anything-tests")
}
func deleteEnvFile() {
	os.Remove(testEnvFile)
}

package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadDbEnv(t *testing.T) {
	envFilePopulate()
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

func TestConnectString(t *testing.T) {
	envFilePopulate()
	godotenv.Load(testEnvFile)
	Vals.loadDBEnv()
	test_str := Vals.DB.ConnectString()
	if test_str != "host=localhost port=5433 user=root password=dbsuperuser991 dbname=comm-anything-tests sslmode=disable" {
		t.Errorf("The connection string was not formatted as expected.")
	}
	deleteEnvFile()
}

func TestLoad(t *testing.T) {
	envFilePopulate()
	defer deleteEnvFile()
	Vals.Load(testEnvFile)
}

func TestLoadServerEnv(t *testing.T) {
	envFilePopulate()
	defer deleteEnvFile()
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
}

// ------------- Test Assistance Functions ---------------------

var testEnvFile = "test.env"

func envFileWrite(key string, value string) {
	f, _ := os.OpenFile(testEnvFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	fmt.Fprintln(f, fmt.Sprintf("%s=%s", key, value))
}

func envFilePopulate() {
	envFileWrite("CA_TESTING_MODE", "true")
	envFileWrite("CA_TESTING_MODE", "true")
	envFileWrite("SERVER_LOG_ALL", "true")
	envFileWrite("DB_IMAGE", "postgres:14.5-alpine")
	envFileWrite("DB_CONTAINER_NAME", "923postgres")
	envFileWrite("DB_CONTAINER_PORT", "5432")
	envFileWrite("DB_HOST", "localhost")
	envFileWrite("DB_HOST_PORT", "5433")
	envFileWrite("DB_USER", "root")
	envFileWrite("DB_PASSWORD", "dbsuperuser991")
	envFileWrite("SERVER_PORT", "3000")
	envFileWrite("JWT_KEY", "key1111111")
	envFileWrite("JWT_COOKIE_NAME", "canywauth")
	envFileWrite("DB_DATABASE_NAME", "comm-anything")
	envFileWrite("TEST_DB_DATABASE_NAME", "comm-anything-tests")
}

func deleteEnvFile() {
	os.Remove(testEnvFile)
}

package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadDbEnv(t *testing.T) {
	setAllEnv()
	defer os.Clearenv()
	err := Vals.loadDBEnv()
	if err != nil {
		t.Errorf("No error should be produced when all .env vars are set")
	}
	if Vals.DB.DBname != "comm-anything-tests" {
		t.Errorf("The db name should be comm-anything-tests.")
	}
	if Vals.DB.Host != "localhost" {
		t.Errorf("The db host should be localhost.")
	}
	if Vals.DB.Password != "dbsuperuser991" {
		t.Errorf("The db password should be dbsuperuser991.")
	}
	if Vals.DB.Port != "5433" {
		t.Errorf("The db port should be 5433.")
	}
	if Vals.DB.User != "root" {
		t.Errorf("The db user should be root.")
	}
	os.Clearenv()
	err = Vals.loadDBEnv()
	if err == nil {
		t.Errorf("An error should be produced if DB_HOST is not set in the .env.")
	}
	os.Setenv("DB_HOST", "localhost")
	err = Vals.loadDBEnv()
	if err == nil {
		t.Errorf("An error should be produced if DB_HOST_PORT is not set in the .env.")
	}
	os.Setenv("DB_HOST_PORT", "5433")
	err = Vals.loadDBEnv()
	if err == nil {
		t.Errorf("An error should be produced if DB_USER is not set in the .env.")
	}
	os.Setenv("DB_USER", "root")
	err = Vals.loadDBEnv()
	if err == nil {
		t.Errorf("An error should be produced if DB_PASSWORD is not set in the .env.")
	}
	os.Setenv("DB_PASSWORD", "dbsuperuser991")
	err = Vals.loadDBEnv()
	if err == nil {
		t.Errorf("An error should be produced if CA_TESTING_MODE is not set in the .env.")
	}
	os.Setenv("CA_TESTING_MODE", "true")
	err = Vals.loadDBEnv()
	if err == nil {
		t.Errorf("An error should be produced if TEST_DB_DATABASE_NAME is not set in the .env and testing mode is on.")
	}
	os.Setenv("TEST_DB_DATABASE_NAME", "comm-anything-tests")
	err = Vals.loadDBEnv()
	if err != nil {
		t.Errorf("No error should be returned.")
	}
	os.Setenv("CA_TESTING_MODE", "false")
	err = Vals.loadDBEnv()
	if err == nil {
		t.Errorf("An error should be produced if DB_DATABASE_NAME is not set in the .env and testing mode is off.")
	}
	os.Setenv("DB_DATABASE_NAME", "comm-anything")
	err = Vals.loadDBEnv()
	if err != nil {
		t.Errorf("No error should be returned.")
	}
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
	err := Vals.Load(testEnvFile)
	if err != nil {
		t.Errorf("No error should result from loading an env file: %s", err.Error())
	}
	deleteEnvFile()
	os.Clearenv()
	err = Vals.Load("...")
	if err == nil {
		t.Errorf("An error should result from a bad env file path.")
	}
	envFileWrite("EXISTS", "true")
	err = Vals.Load(testEnvFile)
	if err == nil {
		t.Errorf("An error should result when an env file exists but isnt populated correctly.")
	}
	os.Clearenv()
	deleteEnvFile()
	envFileDBPopulate()
	err = Vals.Load(testEnvFile)
	if err == nil {
		t.Errorf("An error should result when db stuff is set but server isnt.")
	}
	os.Clearenv()
	deleteEnvFile()

}

func TestReset(t *testing.T) {
	envFilePopulate()
	_ = Vals.Load(testEnvFile)
	Vals.Reset()
	actual := os.Getenv("DB_NAME")
	if actual != "" {
		t.Errorf("Expected env vars to be reset, but got %v", actual)
	}
	if Vals.IsLoaded != false {
		t.Errorf("Expected 'isLoaded' to be set to false in config after reset")
	}
}

func TestLoadServerEnv(t *testing.T) {
	setAllEnv()
	defer os.Clearenv()
	err := Vals.loadServerEnv()
	if err != nil {
		t.Errorf("No error should be produced when all .env vars are set")
	}
	if Vals.Server.DoesLogAll != true {
		t.Errorf("SERVER_LOG_ALL is not correct: %v should be true.", Vals.Server.DoesLogAll)
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
	os.Setenv("SERVER_LOG_ALL", "false")
	err = Vals.loadServerEnv()
	if Vals.Server.DoesLogAll != false {
		t.Errorf("SERVER_LOG_ALL is not correct: %v should be false.", Vals.Server.DoesLogAll)
	}
	os.Clearenv()
	err = Vals.loadServerEnv()
	if err == nil {
		t.Errorf("An error should be produced if SERVER_PORT is not set in the .env.")
	}
	os.Setenv("SERVER_PORT", "3000")
	err = Vals.loadServerEnv()
	if err == nil {
		t.Errorf("An error should be produced if SERVER_LOG_ALL is not set in the .env.")
	}
	os.Setenv("SERVER_LOG_ALL", "true")
	err = Vals.loadServerEnv()
	if err == nil {
		t.Errorf("An error should be produced if JWT_KEY is not set in the .env.")
	}
	os.Setenv("JWT_KEY", "xxxxxx")
	err = Vals.loadServerEnv()
	if err == nil {
		t.Errorf("An error should be produced if JWT_COOKIE_NAME is not set in the .env.")
	}
}

// ------------- Test Assistance Functions ---------------------

var testEnvFile = "test.env"

func envFileWrite(key string, value string) {
	f, _ := os.OpenFile(testEnvFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	fmt.Fprintln(f, fmt.Sprintf("%s=%s", key, value))
}

func setAllEnv() {
	os.Setenv("CA_TESTING_MODE", "true")
	setDBEnv()
	setServEnv()
}

func setDBEnv() {
	os.Setenv("DB_IMAGE", "postgres:14.5-alpine")
	os.Setenv("DB_CONTAINER_NAME", "923postgres")
	os.Setenv("DB_CONTAINER_PORT", "5432")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_HOST_PORT", "5433")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "dbsuperuser991")
	os.Setenv("DB_DATABASE_NAME", "comm-anything")
	os.Setenv("TEST_DB_DATABASE_NAME", "comm-anything-tests")
}
func setServEnv() {
	os.Setenv("SERVER_LOG_ALL", "true")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("JWT_KEY", "key1111111")
	os.Setenv("JWT_COOKIE_NAME", "canywauth")
}

func envFilePopulate() {
	envFileWrite("DB_IMAGE", "postgres:14.5-alpine")
	envFileWrite("DB_CONTAINER_NAME", "923postgres")
	envFileWrite("DB_CONTAINER_PORT", "5432")
	envFileDBPopulate()
	envFileServerPopulate()
}

func envFileDBPopulate() {
	envFileWrite("CA_TESTING_MODE", "true")
	envFileWrite("DB_HOST", "localhost")
	envFileWrite("DB_HOST_PORT", "5433")
	envFileWrite("DB_USER", "root")
	envFileWrite("DB_PASSWORD", "dbsuperuser991")
	envFileWrite("DB_DATABASE_NAME", "comm-anything")
	envFileWrite("TEST_DB_DATABASE_NAME", "comm-anything-tests")

}
func envFileServerPopulate() {
	envFileWrite("SERVER_LOG_ALL", "true")
	envFileWrite("SERVER_PORT", "3000")
	envFileWrite("JWT_KEY", "key1111111")
	envFileWrite("JWT_COOKIE_NAME", "canywauth")

}

func deleteEnvFile() {
	os.Remove(testEnvFile)
}

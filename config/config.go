package config

/*
Config holds values parsed from the .env file. It is used across the application to configure connections. It is accessed through the global singleton config.Vals.
*/
import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DbCredentials are stored in the global Config singleton as Config.DB. It holds the connection settings for accessing the Postgres database.
type DbCredentials struct {
	/* The address the database lives at, e.g., localhost or a url. */
	Host string
	/* The name of the database. */
	DBname string
	/* The port the database is served on. */
	Port string
	/* The username credential. */
	User string
	/* The password credential. */
	Password string
}

// ServerConfig is stored in the global Config singleton as Config.server. It holds the connection settings for the server.
type ServerConfig struct {
	/* The port to serve the server on. */
	Port string
	/* Whether the server will log all incoming requests. */
	DoesLogAll bool
	/* A key to encrypt access tokens (JWT standard) */
	JWTKey string
	/* The cookie name to store the access tokens as on user devices. */
	JWTCookieName string
}

type config struct {
	DB       DbCredentials
	Server   ServerConfig
	IsLoaded bool
}

// Vals is a global configuration object singleton holding environment variables and other global data.
var Vals config

// Load loads the environment variables from the .env file. It should be called in the main function and then in the TestMain function of every package that needs access to those environment variables. While main calls the function with a path to the current working directory, tests will have to use relative directories to find the .env file.
func (c *config) Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	err = Vals.loadDBEnv()
	if err != nil {
		return err
	}
	err = Vals.loadServerEnv()
	if err != nil {
		return err
	}
	Vals.IsLoaded = true
	return nil
}

// loadDBEnv loads database related environment variables into the configuration struct. If it fails to load a variable, it terminates the program process. Correct environment variables are required for the program to run.
func (c *config) loadDBEnv() error {
	c.DB.Host = os.Getenv("DB_HOST")
	if c.DB.Host == "" {
		return getEnvError("DB_HOST")
	}
	c.DB.Port = os.Getenv("DB_HOST_PORT")
	if c.DB.Port == "" {
		return getEnvError("DB_PORT")
	}
	c.DB.User = os.Getenv("DB_USER")
	if c.DB.User == "" {
		return getEnvError("DB_USER")
	}
	c.DB.Password = os.Getenv("DB_PASSWORD")
	if c.DB.Password == "" {
		return getEnvError("DB_PASSWORD")
	}
	testingMode := os.Getenv("CA_TESTING_MODE")
	prodDBname := os.Getenv("DB_DATABASE_NAME")
	testDBname := os.Getenv("TEST_DB_DATABASE_NAME")
	if testingMode == "" || testingMode == "false" || testingMode == "0" {
		if prodDBname == "" {
			return getEnvError("DB_DATABASE_NAME")
		} else {
			c.DB.DBname = prodDBname
		}
	} else if testingMode == "true" || testingMode == "1" {
		if testDBname == "" {
			return getEnvError("TEST_DB_DATABASE_NAME")
		} else {
			c.DB.DBname = testDBname
		}
	} else {
		return getEnvError("CA_TESTING_MODE")
	}
	return nil
}

// DBString builds a string from the database connection credentials and returns it. For use with sql.Open.
func (d *DbCredentials) ConnectString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.DBname)
}

// loadServerEnv loads server related environment variables into the configuration struct. If it fails to load a variable, it terminates the program process. Correct environment variables are required for the program to run.
func (c *config) loadServerEnv() error {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return getEnvError("SERVER_PORT")
	} else {
		Vals.Server.Port = ":" + port
	}
	logall := os.Getenv("SERVER_LOG_ALL")
	if logall == "" {
		return getEnvError("SERVER_LOG_ALL")
	} else if logall == "1" || logall == "true" || logall == "True" {
		c.Server.DoesLogAll = true
	} else {
		c.Server.DoesLogAll = false
	}
	jwtkey := os.Getenv("JWT_KEY")
	if jwtkey == "" {
		return getEnvError("JWT_KEY")
	} else {
		Vals.Server.JWTKey = jwtkey
	}
	cookie_name := os.Getenv("JWT_COOKIE_NAME")
	if cookie_name == "" {
		log.Println("No JWT_COOKIE_NAME in .env. Defaulting to 'canywauth'")
		Vals.Server.JWTCookieName = "canywauth"
	} else {
		Vals.Server.JWTCookieName = cookie_name
	}
	return nil
}

// Gets an error object describing the environmnet variable that was missing or malformed.
func getEnvError(problematic_env_variable_name string) error {
	return errors.New(fmt.Sprintf("Error parsing environment variable %v", problematic_env_variable_name))
}

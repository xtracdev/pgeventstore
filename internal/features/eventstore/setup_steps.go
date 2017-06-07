package eventstore

import (
	. "github.com/gucumber/gucumber"
	"os"
)

var DBUser string
var DBPassword string
var DBHost string
var DBPort string
var DBName string
var configErrors []string

func init() {
	Given(`^some tests to run$`, func() {
	})

	Then(`^the database connection configuration is read from the environment$`, func() {
	})

	GlobalContext.BeforeAll(func() {
		DBUser = os.Getenv("DB_USER")
		if DBUser == "" {
			configErrors = append(configErrors, "Configuration missing DB_USER env variable")
		}

		DBPassword = os.Getenv("DB_PASSWORD")
		if DBPassword == "" {
			configErrors = append(configErrors, "Configuration missing DB_PASSWORD env variable")
		}

		DBHost = os.Getenv("DB_HOST")
		if DBHost == "" {
			configErrors = append(configErrors, "Configuration missing DB_HOST env variable")
		}

		DBPort = os.Getenv("DB_PORT")
		if DBPort == "" {
			configErrors = append(configErrors, "Configuration missing DB_PORT env variable")
		}

		DBName = os.Getenv("DB_NAME")
		if DBName == "" {
			configErrors = append(configErrors, "Configuration missing DB_NAME env variable")
		}

	})

}

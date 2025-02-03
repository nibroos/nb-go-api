package config

import (
	"fmt"
	"os"
)

func GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

func GetTestDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER_TEST"),
		os.Getenv("POSTGRES_PASSWORD_TEST"),
		os.Getenv("POSTGRES_HOST_TEST"),
		os.Getenv("POSTGRES_PORT_TEST"),
		os.Getenv("POSTGRES_DB_TEST"),
	)
}

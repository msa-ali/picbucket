package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	// SMTP Env vars

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	//  DB Env vars

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server Env vars

	ServerUrl  string
	ServerPort string
	CSRFKey    string
	CSRFSecure bool
	ImagesDir  string
}

var env Env

func getInvalidEnvError(prefix string) error {
	return fmt.Errorf("%s: invalid env vars", prefix)
}

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("loadenv: %w", err)
	}
	err = loadSMTPEnv()
	if err != nil {
		return err
	}

	err = loadDBEnv()
	if err != nil {
		return err
	}

	err = loadServerEnv()
	if err != nil {
		return err
	}

	return nil
}

func GetEnv() Env {
	return env
}

func loadSMTPEnv() error {
	env.SMTPHost = os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}
	env.SMTPPort = smtpPort
	env.SMTPUsername = os.Getenv("SMTP_USERNAME")
	env.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	if env.SMTPHost == "" ||
		env.SMTPUsername == "" ||
		env.SMTPPassword == "" {
		return getInvalidEnvError("loadSMTPEnv")
	}
	return nil
}

func loadDBEnv() error {
	env.DBHost = os.Getenv("DB_HOST")
	env.DBPort = os.Getenv("DB_PORT")
	env.DBUser = os.Getenv("DB_USER")
	env.DBPassword = os.Getenv("DB_PASSWORD")
	env.DBName = os.Getenv("DB_NAME")
	env.DBSSLMode = os.Getenv("DB_SSLMODE")

	if env.DBHost == "" ||
		env.DBName == "" ||
		env.DBPassword == "" ||
		env.DBPort == "" ||
		env.DBSSLMode == "" ||
		env.DBUser == "" {
		return getInvalidEnvError("loadDBEnv")
	}
	return nil
}

func loadServerEnv() error {
	env.ServerUrl = os.Getenv("SERVER_URL")
	env.ServerPort = os.Getenv("SERVER_PORT")
	env.CSRFKey = os.Getenv("CSRF_KEY")
	csrfSecure, err := strconv.ParseBool(os.Getenv("CSRF_SECURE"))
	if err != nil {
		return getInvalidEnvError("loadServerEnv")
	}
	env.CSRFSecure = csrfSecure
	env.ImagesDir = os.Getenv("IMAGES_DIR")
	if env.ImagesDir == "" {
		env.ImagesDir = "images"
	}
	if env.ServerUrl == "" ||
		env.ServerPort == "" ||
		env.CSRFKey == "" {
		return getInvalidEnvError("loadServerEnv")
	}
	return nil
}

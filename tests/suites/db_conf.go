package suites

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewDBConf(files ...string) DBConf {
	var err error

	absPaths := make([]string, 0, cap(files))
	for _, f := range files {
		envPath, err := filepath.Abs(f)
		if err != nil {
			panic(err)
		}
		absPaths = append(absPaths, envPath)
	}

	err = godotenv.Load(absPaths...)
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return DBConf{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

package conn

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func InitDB() error {
	conf, err := config()
	if err != nil {
		return err
	}
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			conf["DB_HOST"],
			conf["DB_PORT"],
			conf["DB_USER"],
			conf["DB_PASS"],
			conf["DB_NAME"],
		))
	DB = db
	return nil
}

func config() (map[string]string, error) {
	envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME"}
	conf := make(map[string]string)
	for _, envVar := range envVars {
		value, ok := os.LookupEnv(envVar)
		if !ok {
			return conf, errors.New(fmt.Sprintf("%s enviroment variable required but not set", envVar))
		}
		conf[envVar] = value
	}
	return conf, nil
}

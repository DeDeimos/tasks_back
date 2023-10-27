package dsn

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	DBHost string `json:"DBHost"`
	DBPort string `json:"DBPort"`
	DBUSer string `json:"DBUSer"`
	DBPass string `json:"DBPass"`
	DBName string `json:"DBName"`
}

func SetConnectionString() string {
	file, err := os.Open("./iternal/app/dsn/db_config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return fmt.Sprintf("")
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config JSON:", err)
		return fmt.Sprintf("")
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUSer, config.DBPass, config.DBName)
}

// FromEnv собирает DSN строку из переменных окружения
func FromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

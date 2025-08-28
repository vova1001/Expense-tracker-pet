package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	if _, err := os.Stat("../.env"); err == nil {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env:", err)
		}
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	ssl_mode := os.Getenv("SSL_MODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, db_password, db_name, ssl_mode)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("error connecting to DB", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("ping failed", err)
	}
	fmt.Println("Connection to DB successful")
	_, err = DB.Exec("SET TIME ZONE 'Europe/Moscow'")
	if err != nil {
		log.Fatal("Не удалось установить часовой пояс: ", err)
	}

}

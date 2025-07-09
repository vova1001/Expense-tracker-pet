package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Erorr loading .env")
	}

	host := os.Getenv("DB_HOST")
	fmt.Println(host)

}

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

	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("error connecting to DB", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("ping failed", err)
	}
	fmt.Println("Connection to DB successful")
	defer DB.Close()

	cwd, _ := os.Getwd()
	fmt.Println("Current directory:", cwd)

	BytSQL_create, err := os.ReadFile("../internal/storage/create_and_queries/Create_table.sql")
	if err != nil {
		log.Fatal("Error reading SQL file", err)
	}
	_, err = DB.Exec(string(BytSQL_create))
	if err != nil {
		log.Fatal("Runtime error", err)
	}
	fmt.Println("Table has been successfully filled and created")

	BytSQL_quare, err := os.ReadFile("../internal/storage/create_and_queries/Queries_data.sql")
	if err != nil {
		log.Fatal("Error reading SQL file", err)
	}
	rows, err := DB.Query(string(BytSQL_quare))
	if err != nil {
		log.Fatal("Runtime error qur", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var fio, phone string
		var salary float64

		err = rows.Scan(&id, &fio, &phone, &salary)
		if err != nil {
			log.Fatal("Error reading row", err)
		}
		fmt.Printf("ID=%d, FIO=%s, PHONE=%s, SALARY=%f\n", id, fio, phone, salary)
	}

}

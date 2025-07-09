package storage

import (
	"database/sql"
)

var DB *sql.DB

func InitDB() {
	err := godoenv
}

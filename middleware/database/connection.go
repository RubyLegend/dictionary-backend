package database

import (
	// "database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var db *sql.DB

var cfg = mysql.Config{
	User:                 "",
	Passwd:               "",
	Net:                  "tcp",
	Addr:                 "",
	DBName:               "",
	ParseTime:            true,
	AllowNativePasswords: true,
}

func OpenConnection() bool {
	// If connection is already exists
	if db != nil {
		return true
	}

	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Addr = string(os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"))
	cfg.DBName = os.Getenv("DB_NAME")

	// Open connection
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return false
	}

	// Verifying it
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return false
	}
	log.Println("Database connected!")

	return true
}

func GetConnection() *sql.DB {
	return db
}

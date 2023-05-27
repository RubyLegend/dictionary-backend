package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

var cfg = mysql.Config{
	User:                 "root",
	Passwd:               "YAELfvk5Jgt8qRTc",
	Net:                  "tcp",
	Addr:                 "127.0.0.1:3306",
	DBName:               "Dictionary",
	ParseTime:            true,
	AllowNativePasswords: true,
}

func OpenConnection() bool {
	// If connection is already exists
	if db != nil {
		return true
	}

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

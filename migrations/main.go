package main

import (
	"database/sql"
	"log"
	"os"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	migrate "github.com/golang-migrate/migrate/v4"
	mysql2 "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if ok := OpenConnection(); !ok {
		log.Fatal("Failed to open connection to database.\n",
			"Please check if your server is up and running.")
	}

	driver, _ := mysql2.WithInstance(db, &mysql2.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"mysql", driver)

	if err != nil {
		log.Fatal("Failed to open new migration instance: ", err)
	}

	err = m.Up()

	// Probably due to testing all files in directory this won't work immediately
	if err != nil {
		log.Println("Migration failed: ", err)
	}
}

package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var dbErr error
	DB, dbErr = sql.Open("mysql", dsn)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	fmt.Println("Database berhasil terhubung!")
}

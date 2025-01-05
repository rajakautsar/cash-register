package internal

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

var db *sql.DB

func InitDB() error {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
        return err
    }

    username := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    dbname := os.Getenv("DB_NAME")

    dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)

    db, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }

    fmt.Println("Successfully connected to the database!")
    return nil
}

func GetDB() *sql.DB {
    return db
}
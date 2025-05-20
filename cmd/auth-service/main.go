package main

import (

    "log"
	"fmt"
	"database/sql"
	"github.com/lib/pq"
    //"os"
	"github.com/sashakosti/auth-service/internal/db"
    "github.com/joho/godotenv"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  .env файл не загружен, продолжаем без него")
    }
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL не установлен")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ БД недоступна:", err)
	}

	fmt.Println("✅ Успешное подключение к PostgreSQL!")
}
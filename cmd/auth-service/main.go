package main

import (

    "log"
	"fmt"
	"database/sql"
	"github.com/lib/pq"
    "os"
	"time"
	"github.com/google/uuid"
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

func token_test() {
	
	// Тестовая вставка
	token := db.RefreshToken{
		UserID:    "123e4567-e89b-12d3-a456-426614174000",
		TokenHash: "s3cret",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = db.InsertRefreshToken(conn, token)
	if err != nil {
		log.Fatalf("❌ Ошибка вставки токена: %v", err)
	}

	fmt.Println("✅ Токен вставлен")

	// Тестовый выбор
	rt, err := db.GetRefreshTokenByHash(conn, "s3cret")
	if err != nil {
		log.Fatalf("❌ Ошибка получения токена: %v", err)
	}

	fmt.Printf("🔍 Получен токен: %+v\n", rt)

	// Тестовый удаление
	err = db.DeleteRefreshToken(conn, "s3cret")
	if err != nil {
		log.Fatalf("❌ Ошибка удаления токена: %v", err)
	}

	fmt.Println("✅ Токен удален")
}

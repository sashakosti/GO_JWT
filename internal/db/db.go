package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL не установлен")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("❌ БД не отвечает: %v", err)
	}

	Conn = conn
	fmt.Println("✅ Успешное подключение к PostgreSQL")
}

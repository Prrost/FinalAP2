package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/Prrost/FinalAP2/order-service/infra/logger"
	"github.com/Prrost/FinalAP2/order-service/repository/sqlite"
	grpcsrv "github.com/Prrost/FinalAP2/order-service/transport/grpc"
	"github.com/Prrost/FinalAP2/order-service/usecase"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	grpcPort := getEnv("GRPC_PORT", ":50052")
	dbPath := getEnv("SQLITE_DB_PATH", "orders.db")

	// Открываем SQLite и таблицу
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer db.Close()
	db.Exec(`CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_email TEXT,
        book_id INTEGER,
        taken_at DATETIME,
        due_at DATETIME,
        returned_at DATETIME NULL,
        waiting BOOLEAN
    );`)

	// Репозиторий и юзкейс
	repo := sqlite.NewOrderRepo(db)
	uc := usecase.NewOrderUC(repo)

	// Запуск gRPC
	go func() {
		if err := grpcsrv.RunGRPC(uc, grpcPort); err != nil {
			logger.Log.Fatal(err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logger.Log.Info("Shutting down")
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

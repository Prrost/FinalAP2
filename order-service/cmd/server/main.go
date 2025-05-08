package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/Prrost/FinalAP2/order-service/infra/logger"
	natspub "github.com/Prrost/FinalAP2/order-service/infra/nats"
	"github.com/Prrost/FinalAP2/order-service/repository/sqlite"
	grpcsrv "github.com/Prrost/FinalAP2/order-service/transport/grpc"
	"github.com/Prrost/FinalAP2/order-service/usecase"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nats-io/nats.go"
)

func main() {
	grpcPort := getEnv("GRPC_PORT", ":50052")
	natsURL := getEnv("NATS_URL", nats.DefaultURL)
	dbPath := getEnv("SQLITE_DB_PATH", "orders.db")

	// Открываем SQLite и таблицу
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer db.Close()
	db.Exec(`CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        book_id INTEGER,
        taken_at DATETIME,
        due_at DATETIME,
        returned_at DATETIME NULL,
        waiting BOOLEAN
    );`)

	// NATS publisher
	pub, err := natspub.NewPublisher(natsURL)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// Репозиторий и юзкейс
	repo := sqlite.NewOrderRepo(db)
	uc := usecase.NewOrderUC(repo, pub)

	// Подписка на book.available
	nc, _ := nats.Connect(natsURL)
	nc.Subscribe("book.available", func(msg *nats.Msg) {
		var evt struct{ BookID int64 }
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			logger.Log.Errorf("book.available: %v", err)
			return
		}
		waiting, _ := repo.ListWaiting(evt.BookID)
		for _, o := range waiting {
			logger.Log.Infof("Notify user %d: book %d available", o.UserID, evt.BookID)
		}
	})

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

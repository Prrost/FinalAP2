// cmd/server/main.go
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/Prrost/FinalAP2/order-service/httpserver"
	"github.com/Prrost/FinalAP2/order-service/infra/email"
	"github.com/Prrost/FinalAP2/order-service/infra/logger"
	sqliteRepo "github.com/Prrost/FinalAP2/order-service/repository/sqlite"
	grpcsrv "github.com/Prrost/FinalAP2/order-service/transport/grpc"
	"github.com/Prrost/FinalAP2/order-service/usecase"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// флаги для SMTP
	smtpHost := flag.String("smtp-host", "", "SMTP host (e.g. smtp.gmail.com)")
	smtpPort := flag.String("smtp-port", "587", "SMTP port")
	smtpUser := flag.String("smtp-user", "", "SMTP username/email")
	smtpPass := flag.String("smtp-pass", "", "SMTP app password")
	grpcPort := flag.String("grpc-port", ":50052", "gRPC listen address")
	httpPort := flag.String("http-port", ":8080", "HTTP listen address")
	dbPath := flag.String("db", "orders.db", "SQLite database path")
	flag.Parse()

	// опционально загружаем .env, чтобы иметь запасной вариант
	_ = godotenv.Load()

	// Проверяем, что флаги заданы
	if *smtpHost == "" || *smtpUser == "" || *smtpPass == "" {
		log.Fatal("you must provide -smtp-host, -smtp-user and -smtp-pass")
	}

	// Открываем БД
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer db.Close()

	// Создаём таблицу
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_email TEXT,
            book_id INTEGER,
            taken_at DATETIME,
            due_at DATETIME,
            returned_at DATETIME NULL,
            waiting BOOLEAN
        );
    `)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// Репозиторий
	repo := sqliteRepo.NewOrderRepo(db)

	// SMTP-клиент из флагов
	emailSender := email.NewEmailSenderConfig(
		*smtpHost,
		*smtpPort,
		*smtpUser,
		*smtpPass,
	)
	// Мини-тест отправки сразу при старте
	log.Printf("[SMTP TEST] host=%q user=%q passSet=%t",
		emailSender.Host, emailSender.Username, emailSender.Password != "",
	)
	if err := emailSender.Send(
		emailSender.Username,
		"SMTP Test",
		"Если вы это письмо получили — SMTP работает!",
	); err != nil {
		log.Fatalf("SMTP тест не прошёл: %v", err)
	}
	log.Println("SMTP тест прошёл — письмо отправлено!")

	// Usecase
	uc, err := usecase.NewOrderUC(repo, emailSender)
	if err != nil {
		logger.Log.Fatal(err)
	}

	// Запускаем HTTP
	go func() {
		logger.Log.Infof("HTTP listening on %s", *httpPort)
		httpserver.RunHTTP(uc, *httpPort)
	}()
	// Запускаем gRPC
	go func() {
		logger.Log.Infof("gRPC listening on %s", *grpcPort)
		if err := grpcsrv.RunGRPC(uc, *grpcPort); err != nil {
			logger.Log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down")
}

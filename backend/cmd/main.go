package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"movie-quizzer/backend/internal/config"
	"movie-quizzer/backend/internal/server"
	"movie-quizzer/backend/internal/service"
	"movie-quizzer/backend/internal/storage"
	"movie-quizzer/backend/internal/storage/filestorage"
)

func main() {
	cfg := config.Load()

	// Определяем sslmode в зависимости от окружения
	sslmode := "disable"
	if cfg.DBHost != "postgres" && cfg.DBHost != "localhost" {
		// Для внешних БД (например, Render) используем require
		sslmode = "require"
	}
	
	dsn := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=" + sslmode

	// Retry подключения к БД
	var db *sqlx.DB
	maxRetries := 30
	connected := false
	for i := 0; i < maxRetries; i++ {
		var err error
		db, err = sqlx.Open("pgx", dsn)
		if err != nil {
			log.Printf("db open attempt %d/%d failed: %v", i+1, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)

		if err := db.Ping(); err != nil {
			log.Printf("db ping attempt %d/%d failed: %v", i+1, maxRetries, err)
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		connected = true
		break
	}
	if !connected {
		log.Fatalf("failed to connect to database after %d attempts", maxRetries)
	}
	log.Println("database connection established")

	dbStore := storage.NewSQL(db)

	if err := dbStore.RunMigrations(); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	// Retry подключения к MinIO
	var minioClient *filestorage.MinIO
	maxMinioRetries := 10
	minioConnected := false
	for i := 0; i < maxMinioRetries; i++ {
		var err error
		minioClient, err = filestorage.New(
			cfg.MinioEndpoint,
			cfg.MinioAccessKey,
			cfg.MinioSecretKey,
			cfg.MinioBucket,
			cfg.MinioUseSSL,
		)
		if err != nil {
			log.Printf("minio connection attempt %d/%d failed: %v", i+1, maxMinioRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}
		minioConnected = true
		break
	}
	if !minioConnected {
		log.Fatalf("failed to connect to minio after %d attempts", maxMinioRetries)
	}
	log.Println("minio connection established")

	svc := service.New(dbStore, minioClient)

	s := server.New(svc)
	router := s.Router()

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatalf("http server: %v", err)
	}
}

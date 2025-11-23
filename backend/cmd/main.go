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

	dsn := "host=" + cfg.DBHost +
		" port=" + cfg.DBPort +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=disable"

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	dbStore := storage.NewSQL(db)

	if err := dbStore.RunMigrations(); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	minioClient, err := filestorage.New(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		false,
	)
	if err != nil {
		log.Fatalf("minio: %v", err)
	}

	svc := service.New(dbStore, minioClient)

	s := server.New(svc)
	router := s.Router()

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatalf("http server: %v", err)
	}
}

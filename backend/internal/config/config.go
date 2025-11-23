package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	HTTPAddr   string `env:"HTTP_ADDR"       env-default:":8080"`
	DBHost     string `env:"DB_HOST"         env-default:"postgres"`
	DBPort     string `env:"DB_PORT"         env-default:"5432"`
	DBUser     string `env:"DB_USER"         env-default:"app"`
	DBPassword string `env:"DB_PASSWORD"     env-default:"app"`
	DBName     string `env:"DB_NAME"         env-default:"quiz"`

	MinioEndpoint  string `env:"MINIO_ENDPOINT"  env-default:"minio:9000"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY" env-default:"minio"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY" env-default:"minio123"`
	MinioBucket    string `env:"MINIO_BUCKET"    env-default:"media"`
}

func Load() Config {
	var c Config
	_ = cleanenv.ReadEnv(&c)
	return c
}

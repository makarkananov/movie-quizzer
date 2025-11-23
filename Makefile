.PHONY: help build up down restart logs clean seed

help:
	@echo "Доступные команды:"
	@echo "  make build      - Собрать все Docker образы"
	@echo "  make up         - Запустить все сервисы"
	@echo "  make down       - Остановить все сервисы"
	@echo "  make restart    - Перезапустить все сервисы"
	@echo "  make logs       - Показать логи всех сервисов"
	@echo "  make clean      - Удалить все контейнеры и volumes"
	@echo "  make seed       - Добавить тестовые данные в БД"
	@echo "  make upload-media - Загрузить медиафайлы в MinIO"
	@echo "  make seed-all   - Добавить данные и загрузить медиа"

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	docker system prune -f

seed:
	@echo "Добавление тестовых данных..."
	@docker exec -i quiz-postgres psql -U app -d quiz < backend/scripts/seed.sql || echo "Ошибка: убедитесь, что контейнер quiz-postgres запущен"

upload-media:
	@echo "Загрузка медиафайлов в MinIO..."
	@if [ -f scripts/upload-media-docker.sh ]; then \
		chmod +x scripts/upload-media-docker.sh && \
		./scripts/upload-media-docker.sh; \
	else \
		echo "Используйте MinIO Console: http://localhost:9090 (minio/minio123)"; \
	fi

seed-all: seed upload-media
	@echo "Готово! Данные и медиафайлы загружены."

check-media:
	@echo "Проверка файлов в MinIO..."
	@if [ -f scripts/check-media.ps1 ]; then \
		powershell -ExecutionPolicy Bypass -File scripts/check-media.ps1; \
	elif [ -f scripts/check-media.sh ]; then \
		chmod +x scripts/check-media.sh && ./scripts/check-media.sh; \
	else \
		echo "Используйте MinIO Console: http://localhost:9090"; \
	fi


#!/bin/bash

# Скрипт для проверки загруженных файлов в MinIO
# Использование: ./check-media.sh

BUCKET="media"
NETWORK_NAME="movie-quizzer_default"

echo "Проверка файлов в MinIO..."

# Настройка MinIO Client
docker run --rm --network $NETWORK_NAME \
    minio/mc alias set localminio http://minio:9000 minio minio123 2>/dev/null

echo ""
echo "Файлы в frames/:"
docker run --rm --network $NETWORK_NAME \
    minio/mc ls localminio/$BUCKET/frames/ 2>/dev/null || echo "Папка frames/ пуста или не найдена"

echo ""
echo "Файлы в videos/:"
docker run --rm --network $NETWORK_NAME \
    minio/mc ls localminio/$BUCKET/videos/ 2>/dev/null || echo "Папка videos/ пуста или не найдена"

echo ""
echo "Проверка завершена!"


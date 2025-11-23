#!/bin/bash

# Скрипт для загрузки медиафайлов в MinIO
# Использование: ./upload-media.sh <путь_к_папке_с_медиа>

MINIO_ENDPOINT="localhost:9000"
MINIO_ACCESS_KEY="minio"
MINIO_SECRET_KEY="minio123"
BUCKET="media"

# Проверка наличия mc (MinIO Client)
if ! command -v mc &> /dev/null; then
    echo "MinIO Client (mc) не установлен."
    echo "Установите его с https://min.io/docs/minio/linux/reference/minio-mc.html"
    exit 1
fi

mc alias set localminio http://${MINIO_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} 2>/dev/null || true

mc mb localminio/${BUCKET} 2>/dev/null || true

MEDIA_DIR=${1:-"./media"}

if [ ! -d "$MEDIA_DIR" ]; then
    echo "Папка $MEDIA_DIR не найдена!"
    echo "Создайте папку media/ с подпапками frames/ и videos/"
    exit 1
fi

echo "Загрузка медиафайлов из $MEDIA_DIR..."

# Загрузка изображений
if [ -d "$MEDIA_DIR/frames" ]; then
    echo "Загрузка изображений из frames/..."
    mc cp --recursive "$MEDIA_DIR/frames/" localminio/${BUCKET}/frames/
fi

# Загрузка видео
if [ -d "$MEDIA_DIR/videos" ]; then
    echo "Загрузка видео из videos/..."
    mc cp --recursive "$MEDIA_DIR/videos/" localminio/${BUCKET}/videos/
fi

echo "Готово! Медиафайлы загружены в MinIO bucket '${BUCKET}'"


#!/bin/bash

# Скрипт для загрузки медиафайлов в MinIO через Docker
# Использование: ./upload-media-docker.sh

MINIO_ENDPOINT="minio:9000"
MINIO_ACCESS_KEY="minio"
MINIO_SECRET_KEY="minio123"
BUCKET="media"
MEDIA_DIR="./media"

echo "Загрузка медиафайлов в MinIO через Docker..."

# Проверка наличия папки media
if [ ! -d "$MEDIA_DIR" ]; then
    echo "Папка $MEDIA_DIR не найдена!"
    echo "Создайте папку media/ с подпапками frames/ и videos/"
    exit 1
fi

# Настройка alias через Docker
docker run --rm --network movie-quizzer_default \
    minio/mc alias set localminio http://${MINIO_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} 2>/dev/null || true

# Создание bucket
docker run --rm --network movie-quizzer_default \
    minio/mc mb localminio/${BUCKET} 2>/dev/null || true

# Загрузка изображений (исключаем .gitkeep)
if [ -d "$MEDIA_DIR/frames" ]; then
    echo "Загрузка изображений из frames/..."
    # Используем один контейнер для всех операций с правильным entrypoint
    docker run --rm -v "$(pwd)/media:/data" \
        --network movie-quizzer_default \
        --entrypoint /bin/sh \
        minio/mc -c "
            mc alias set localminio http://${MINIO_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY}
            for file in /data/frames/*; do
                if [ -f \"\$file\" ] && [ \"\$(basename \"\$file\")\" != \".gitkeep\" ]; then
                    filename=\$(basename \"\$file\")
                    echo \"  Загрузка: \$filename\"
                    mc cp \"\$file\" localminio/${BUCKET}/frames/\$filename
                fi
            done
        "
fi

# Загрузка видео (исключаем .gitkeep)
if [ -d "$MEDIA_DIR/videos" ]; then
    echo "Загрузка видео из videos/..."
    # Используем один контейнер для всех операций с правильным entrypoint
    docker run --rm -v "$(pwd)/media:/data" \
        --network movie-quizzer_default \
        --entrypoint /bin/sh \
        minio/mc -c "
            mc alias set localminio http://${MINIO_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY}
            for file in /data/videos/*; do
                if [ -f \"\$file\" ] && [ \"\$(basename \"\$file\")\" != \".gitkeep\" ]; then
                    filename=\$(basename \"\$file\")
                    echo \"  Загрузка: \$filename\"
                    mc cp \"\$file\" localminio/${BUCKET}/videos/\$filename
                fi
            done
        "
fi

echo ""
echo "Проверка загруженных файлов..."
echo "Файлы в frames/:"
docker run --rm --network movie-quizzer_default \
    --entrypoint /bin/sh \
    minio/mc -c "mc alias set localminio http://${MINIO_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} && mc ls localminio/${BUCKET}/frames/" || echo "Папка frames/ пуста или не найдена"
echo ""
echo "Файлы в videos/:"
docker run --rm --network movie-quizzer_default \
    --entrypoint /bin/sh \
    minio/mc -c "mc alias set localminio http://${MINIO_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} && mc ls localminio/${BUCKET}/videos/" || echo "Папка videos/ пуста или не найдена"

echo ""
echo "Готово! Медиафайлы загружены в MinIO bucket '${BUCKET}'"
echo "Проверьте в MinIO Console: http://localhost:9090"


#!/bin/sh

echo "=== Загрузка медиафайлов в MinIO ==="

# Ждем, пока MinIO будет готов
echo "Ожидание готовности MinIO..."
for i in 1 2 3 4 5 6 7 8 9 10; do
  if mc alias set localminio http://minio:9000 minio minio123 2>/dev/null; then
    # Проверяем, что можем подключиться
    if mc ls localminio 2>/dev/null >/dev/null; then
      echo "MinIO готов и доступен!"
      break
    fi
  fi
  if [ $i -eq 10 ]; then
    echo "ОШИБКА: Не удалось подключиться к MinIO после 10 попыток"
    exit 1
  fi
  echo "Попытка $i/10..."
  sleep 2
done

# Проверяем, что bucket существует
if ! mc ls localminio/media 2>/dev/null >/dev/null; then
  echo "Создание bucket media..."
  mc mb localminio/media 2>/dev/null || true
fi

echo ""

# Загрузка изображений
total_uploaded=0
if [ -d /data/frames ]; then
  echo "Загрузка изображений из frames/..."
  count=0
  uploaded=0
  for file in /data/frames/*; do
    if [ -f "$file" ] && [ "$(basename "$file")" != ".gitkeep" ]; then
      filename=$(basename "$file")
      count=$((count + 1))
      echo "  [$count] Загрузка: $filename"
      if mc cp "$file" localminio/media/frames/$filename 2>&1; then
        uploaded=$((uploaded + 1))
        total_uploaded=$((total_uploaded + 1))
      else
        echo "    ОШИБКА загрузки $filename"
      fi
    fi
  done
  echo "  Загружено изображений: $uploaded/$count"
else
  echo "Папка frames/ пуста или не найдена, пропускаем..."
fi

echo ""

# Загрузка видео
if [ -d /data/videos ]; then
  echo "Загрузка видео из videos/..."
  count=0
  uploaded=0
  for file in /data/videos/*; do
    if [ -f "$file" ] && [ "$(basename "$file")" != ".gitkeep" ]; then
      filename=$(basename "$file")
      count=$((count + 1))
      echo "  [$count] Загрузка: $filename"
      if mc cp "$file" localminio/media/videos/$filename 2>&1; then
        uploaded=$((uploaded + 1))
        total_uploaded=$((total_uploaded + 1))
      else
        echo "    ОШИБКА загрузки $filename"
      fi
    fi
  done
  echo "  Загружено видео: $uploaded/$count"
else
  echo "Папка videos/ пуста или не найдена, пропускаем..."
fi

echo ""
if [ $total_uploaded -gt 0 ]; then
  echo "=== Успешно загружено файлов: $total_uploaded ==="
else
  echo "=== Предупреждение: файлы не загружены (папки пусты или файлы отсутствуют) ==="
fi


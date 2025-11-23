# PowerShell скрипт для загрузки медиафайлов в MinIO через Docker
# Использование: .\upload-media.ps1

$MINIO_ENDPOINT = "localhost:9000"
$MINIO_ACCESS_KEY = "minio"
$MINIO_SECRET_KEY = "minio123"
$BUCKET = "media"
$MEDIA_DIR = ".\media"

Write-Host "Загрузка медиафайлов в MinIO..." -ForegroundColor Green

# Проверка наличия папки media
if (-not (Test-Path $MEDIA_DIR)) {
    Write-Host "Папка $MEDIA_DIR не найдена!" -ForegroundColor Red
    Write-Host "Создайте папку media/ с подпапками frames/ и videos/" -ForegroundColor Yellow
    exit 1
}

# Получаем имя сети Docker Compose (пробуем разные варианты)
$NETWORK_NAME = docker network ls --filter "name=movie-quizzer" --format "{{.Name}}" | Select-Object -First 1
if (-not $NETWORK_NAME) {
    $NETWORK_NAME = docker network ls --filter "name=quiz" --format "{{.Name}}" | Select-Object -First 1
}
if (-not $NETWORK_NAME) {
    $NETWORK_NAME = "movie-quizzer_default"
}

Write-Host "Используется сеть: $NETWORK_NAME" -ForegroundColor Cyan

# Проверяем, запущен ли MinIO
$minioRunning = docker ps --filter "name=minio" --format "{{.Names}}" | Select-Object -First 1
if (-not $minioRunning) {
    Write-Host "Ошибка: MinIO контейнер не запущен!" -ForegroundColor Red
    Write-Host "Запустите: docker compose up -d" -ForegroundColor Yellow
    exit 1
}

# Настройка MinIO Client (в одном контейнере для сохранения alias)
Write-Host "Настройка MinIO Client..." -ForegroundColor Cyan
docker run --rm --network $NETWORK_NAME `
    minio/mc alias set localminio http://minio:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY 2>&1 | Out-Null

# Создание bucket (игнорируем ошибку, если уже существует)
Write-Host "Создание bucket..." -ForegroundColor Cyan
docker run --rm --network $NETWORK_NAME `
    minio/mc mb localminio/$BUCKET 2>&1 | Out-Null

# Загрузка изображений
if (Test-Path "$MEDIA_DIR\frames") {
    Write-Host "Загрузка изображений из frames/..." -ForegroundColor Cyan
    $framesPath = (Resolve-Path "$MEDIA_DIR\frames").Path
    $files = Get-ChildItem "$framesPath" -File
    if ($files.Count -eq 0) {
        Write-Host "  Предупреждение: папка frames/ пуста!" -ForegroundColor Yellow
    } else {
        # Фильтруем файлы, исключая .gitkeep
        $imageFiles = $files | Where-Object { $_.Name -ne '.gitkeep' }
        Write-Host "  Найдено файлов (без .gitkeep): $($imageFiles.Count)" -ForegroundColor Cyan
        foreach ($file in $imageFiles) {
            Write-Host "  Загрузка: $($file.Name)" -ForegroundColor Gray
            docker run --rm -v "${framesPath}:/data/frames" `
                --network $NETWORK_NAME `
                --entrypoint /bin/sh `
                minio/mc -c "mc alias set localminio http://minio:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY && mc cp /data/frames/$($file.Name) localminio/$BUCKET/frames/$($file.Name)"
        }
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✓ Изображения загружены успешно" -ForegroundColor Green
        } else {
            Write-Host "  ✗ Ошибка загрузки изображений" -ForegroundColor Red
        }
    }
} else {
    Write-Host "Папка frames/ не найдена, пропускаем..." -ForegroundColor Yellow
}

# Загрузка видео
if (Test-Path "$MEDIA_DIR\videos") {
    Write-Host "Загрузка видео из videos/..." -ForegroundColor Cyan
    $videosPath = (Resolve-Path "$MEDIA_DIR\videos").Path
    $files = Get-ChildItem "$videosPath" -File
    if ($files.Count -eq 0) {
        Write-Host "  Предупреждение: папка videos/ пуста!" -ForegroundColor Yellow
    } else {
        # Фильтруем файлы, исключая .gitkeep
        $videoFiles = $files | Where-Object { $_.Name -ne '.gitkeep' }
        Write-Host "  Найдено файлов (без .gitkeep): $($videoFiles.Count)" -ForegroundColor Cyan
        foreach ($file in $videoFiles) {
            Write-Host "  Загрузка: $($file.Name)" -ForegroundColor Gray
            docker run --rm -v "${videosPath}:/data/videos" `
                --network $NETWORK_NAME `
                --entrypoint /bin/sh `
                minio/mc -c "mc alias set localminio http://minio:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY && mc cp /data/videos/$($file.Name) localminio/$BUCKET/videos/$($file.Name)"
        }
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✓ Видео загружены успешно" -ForegroundColor Green
        } else {
            Write-Host "  ✗ Ошибка загрузки видео" -ForegroundColor Red
        }
    }
} else {
    Write-Host "Папка videos/ не найдена, пропускаем..." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Проверка загруженных файлов..." -ForegroundColor Cyan
Write-Host "Файлы в frames/:" -ForegroundColor Yellow
docker run --rm --network $NETWORK_NAME `
    --entrypoint /bin/sh `
    minio/mc -c "mc alias set localminio http://minio:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY && mc ls localminio/$BUCKET/frames/" 2>&1
Write-Host ""
Write-Host "Файлы в videos/:" -ForegroundColor Yellow
docker run --rm --network $NETWORK_NAME `
    --entrypoint /bin/sh `
    minio/mc -c "mc alias set localminio http://minio:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY && mc ls localminio/$BUCKET/videos/" 2>&1

Write-Host ""
Write-Host "Готово! Медиафайлы загружены в MinIO bucket '$BUCKET'" -ForegroundColor Green
Write-Host "Проверьте в MinIO Console: http://localhost:9090" -ForegroundColor Cyan


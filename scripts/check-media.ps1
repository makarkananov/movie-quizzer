# Скрипт для проверки загруженных файлов в MinIO
# Использование: .\check-media.ps1

$BUCKET = "media"

Write-Host "Проверка файлов в MinIO..." -ForegroundColor Green

# Получаем имя сети Docker Compose
$NETWORK_NAME = docker network ls --filter "name=movie-quizzer" --format "{{.Name}}" | Select-Object -First 1
if (-not $NETWORK_NAME) {
    $NETWORK_NAME = "movie-quizzer_default"
}

Write-Host "Используется сеть: $NETWORK_NAME" -ForegroundColor Cyan
Write-Host ""

# Настройка MinIO Client
docker run --rm --network $NETWORK_NAME `
    minio/mc alias set localminio http://minio:9000 minio minio123 2>$null

Write-Host "Файлы в frames/:" -ForegroundColor Yellow
docker run --rm --network $NETWORK_NAME `
    minio/mc ls localminio/$BUCKET/frames/ 2>$null

Write-Host ""
Write-Host "Файлы в videos/:" -ForegroundColor Yellow
docker run --rm --network $NETWORK_NAME `
    minio/mc ls localminio/$BUCKET/videos/ 2>$null

Write-Host ""
Write-Host "Проверка завершена!" -ForegroundColor Green


# Инструкция по запуску автоматического тестирования API

## Для Linux/WSL/Mac (Bash)

```bash
# Перейдите в директорию проекта
cd movie-quizzer

# Сделайте скрипт исполняемым (если еще не сделано)
chmod +x scripts/test-api.sh

# Запустите тестирование
./scripts/test-api.sh
```

Или из любой директории:

```bash
bash scripts/test-api.sh
```

## Для Windows (PowerShell)

```powershell
# Перейдите в директорию проекта
cd movie-quizzer

# Запустите тестирование
.\scripts\test-api.ps1
```

## Требования

### Для bash версии:
- `curl` - для HTTP запросов
- `bc` (опционально) - для вычисления процентов

Установка в Ubuntu/Debian:
```bash
sudo apt-get update
sudo apt-get install -y curl bc
```

Установка в WSL:
```bash
sudo apt update
sudo apt install -y curl bc
```

### Для PowerShell версии:
- PowerShell 5.1+ или PowerShell Core

## Проверка перед запуском

Убедитесь, что проект запущен:

```bash
# Проверка health check
curl http://localhost:8080/health

# Должен вернуть: {"status":"ok"}
```

Если проект не запущен:

```bash
# Запуск через Docker Compose
docker-compose up -d

# Или через Make
make up
```

## Пример вывода

```
=== Тестирование API Movie Quizzer ===

=== 1. Проверка Health Check ===
[1] Health Check ... ✓ ПРОЙДЕН

=== 2. Тестирование регистрации ===
[2] Регистрация нового пользователя ... ✓ ПРОЙДЕН
[3] Регистрация с невалидным email ... ✓ ПРОЙДЕН
...

=== Итоги тестирования ===
Всего тестов: 20
Пройдено: 19
Провалено: 1
Процент успешности: 95.0%
```

## Устранение проблем

### Ошибка "command not found: curl"
Установите curl:
```bash
sudo apt-get install curl  # Ubuntu/Debian
brew install curl         # Mac
```

### Ошибка "Connection refused"
Убедитесь, что backend запущен:
```bash
docker-compose ps
# Должен показать запущенные контейнеры
```

### Ошибка "Permission denied"
Сделайте скрипт исполняемым:
```bash
chmod +x scripts/test-api.sh
```

### Ошибка с путями в WSL
Используйте правильные пути:
```bash
# В WSL используйте Linux пути
cd /mnt/c/Users/nikha/GolandProjects/movie-quizzer
./scripts/test-api.sh
```


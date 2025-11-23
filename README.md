# Movie Quizzer

Квиз-игра о фильмах с тремя игровыми режимами: по кадру, по видеофрагменту и по цитате.

## Технологический стек

- **Backend**: Go (стандартная библиотека)
- **Frontend**: Vue.js 3 + TypeScript + Vite
- **База данных**: PostgreSQL
- **Хранилище медиа**: MinIO (S3-совместимое)
- **Инфраструктура**: Docker Compose

## Структура проекта

```
movie-quizzer/
├── backend/          # Go backend сервер
├── frontend/         # Vue.js frontend приложение
├── docker-compose.yml
└── Dockerfile
```

## Запуск проекта

### Требования

- Docker и Docker Compose
- Node.js 20+ (для разработки frontend)

### Запуск через Docker Compose

1. Клонируйте репозиторий
2. Запустите все сервисы с пересборкой:

```bash
docker-compose up --build
```

Или в фоновом режиме:

```bash
docker-compose up --build -d
```

Сервисы будут доступны:
- Backend API: http://localhost:8080
- Frontend: http://localhost:3000
- MinIO Console: http://localhost:9090 (minio/minio123)
- PostgreSQL: localhost:5432 (app/app)

После первого запуска добавьте тестовые данные:

```bash
# Добавить вопросы в БД
make seed

# Загрузить медиафайлы в MinIO
make upload-media

# Или все сразу:
make seed-all
```

**Важно:** Перед загрузкой медиафайлов поместите их в папку `media/`:
- Изображения в `media/frames/` (matrix.jpg, inception.jpg, и т.д.)
- Видео в `media/videos/` (avengers.mp4, interstellar.mp4, и т.д.)

См. `media/FILES_NEEDED.md` для списка необходимых файлов.

Или загрузите вручную через MinIO Console: http://localhost:9090 (minio/minio123)

### Разработка Frontend

Для разработки frontend локально:

```bash
cd frontend
npm install
npm run dev
```

Frontend будет доступен на http://localhost:3000 с проксированием API запросов на backend.

### Сборка Frontend для production

```bash
cd frontend
npm install
npm run build
```

Собранные файлы будут в `frontend/dist/`.

## API Endpoints

### Авторизация
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход
- `GET /api/auth/me` - Текущий пользователь

### Игра
- `POST /api/game/sessions` - Начать игровую сессию
- `POST /api/game/sessions/{id}/answers` - Отправить ответ
- `GET /api/game/sessions/{id}/next` - Следующий вопрос
- `GET /api/game/sessions/{id}` - Итоги сессии

### Профиль
- `GET /api/profile` - Статистика пользователя
- `GET /api/profile/achievements` - Достижения

### Рейтинг
- `GET /api/leaderboard/global` - Глобальный рейтинг
- `GET /api/leaderboard/me` - Позиция пользователя

### Медиа
- `GET /api/media/{bucket}/{file}` - Получить медиафайл

## Переменные окружения

Backend использует следующие переменные окружения (значения по умолчанию):

- `HTTP_ADDR=:8080` - Адрес HTTP сервера
- `DB_HOST=postgres` - Хост PostgreSQL
- `DB_PORT=5432` - Порт PostgreSQL
- `DB_USER=app` - Пользователь БД
- `DB_PASSWORD=app` - Пароль БД
- `DB_NAME=quiz` - Имя БД
- `MINIO_ENDPOINT=minio:9000` - Endpoint MinIO
- `MINIO_ACCESS_KEY=minio` - Access Key MinIO
- `MINIO_SECRET_KEY=minio123` - Secret Key MinIO
- `MINIO_BUCKET=media` - Bucket для медиафайлов

## Игровые режимы

1. **По кадру** (`frame`) - Угадайте фильм по изображению кадра
2. **По видеофрагменту** (`video`) - Угадайте фильм по видеоролику
3. **По цитате** (`quote`) - Угадайте фильм по известной цитате

Каждая игровая сессия состоит из 10 вопросов с таймером 60 секунд на вопрос.

## Система начисления очков

- Базовые 10 очков за правильный ответ
- Бонус до 5 очков за скорость (быстрее 30 секунд)
- 0 очков за неправильный ответ

## Разработка

### Backend

```bash
cd backend
go mod download
go run cmd/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Авторы

- Кананов Макар Кириллович (M3406)
- Хабаров Никита Игоревич (M3405)
- Козлов Андрей Викторович (M3406)
- Шафиков Евгений Наильевич (M3406)

## Лицензия

Учебный проект для ИТМО


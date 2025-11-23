# Быстрый старт

## 1. Запуск проекта

```bash
docker-compose up --build
```

## 2. Добавление вопросов в БД

```bash
make seed
```

Или вручную:
```bash
docker exec -i quiz-postgres psql -U app -d quiz < backend/scripts/seed.sql
```

## 3. Добавление медиафайлов (картинок и видео)

### Шаг 1: Поместите файлы в папку `media/`

Создайте структуру:
```
media/
├── frames/          # Изображения
│   ├── matrix.jpg
│   ├── inception.jpg
│   ├── pulp_fiction.jpg
│   ├── shawshank.jpg
│   └── godfather.jpg
└── videos/          # Видео
    ├── avengers.mp4
    ├── interstellar.mp4
    ├── joker.mp4
    ├── parasite.mp4
    └── dune.mp4
```

**Важно:** Имена файлов должны точно совпадать с путями в `backend/scripts/seed.sql`!

### Шаг 2: Загрузите в MinIO

**Способ 1 - Автоматически (рекомендуется):**

```bash
# Windows PowerShell
.\scripts\upload-media.ps1

# Linux/Mac
make upload-media
```

**Способ 2 - Через веб-интерфейс:**

1. Откройте http://localhost:9090
2. Войдите: `minio` / `minio123`
3. Выберите bucket `media`
4. Нажмите "Upload" и загрузите файлы:
   - Изображения в папку `frames/` (например: `frames/matrix.jpg`)
   - Видео в папку `videos/` (например: `videos/avengers.mp4`)

**Способ 3 - Все сразу:**

```bash
make seed-all  # Загрузит данные и медиафайлы
```

## 4. Проверка работы

1. Откройте http://localhost:3000
2. Зарегистрируйтесь или войдите
3. Выберите режим игры
4. Проверьте, что картинки/видео загружаются

## Решение проблем

### Картинки не отображаются

1. Проверьте, что файлы загружены в MinIO (http://localhost:9090)
2. Убедитесь, что пути в БД совпадают с именами файлов
3. Проверьте консоль браузера (F12) на ошибки 404
4. Проверьте логи backend: `docker-compose logs backend`

### Кнопка "Ответить" не работает

1. Откройте консоль браузера (F12)
2. Проверьте ошибки JavaScript
3. Убедитесь, что вопрос загружен (в консоли должно быть `question.value`)

### Нет вопросов в игре

1. Убедитесь, что выполнили `make seed`
2. Проверьте БД: `docker exec -it quiz-postgres psql -U app -d quiz -c "SELECT COUNT(*) FROM questions;"`

### Файлы не находятся в MinIO

1. Проверьте логи backend: `docker-compose logs backend | grep Media`
2. Убедитесь, что файлы загружены в правильный bucket (`media`)
3. Проверьте пути - они должны быть `frames/matrix.jpg`, а не `media/frames/matrix.jpg`

Подробнее см. [MEDIA_SETUP.md](MEDIA_SETUP.md)

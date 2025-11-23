# Инструкция по добавлению медиафайлов

## Структура папок

Создайте следующую структуру папок в корне проекта:

```
movie-quizzer/
├── media/
│   ├── frames/          # Изображения кадров из фильмов
│   │   ├── matrix.jpg
│   │   ├── inception.jpg
│   │   ├── pulp_fiction.jpg
│   │   ├── shawshank.jpg
│   │   └── godfather.jpg
│   └── videos/          # Видеофрагменты из фильмов
│       ├── avengers.mp4
│       ├── interstellar.mp4
│       ├── joker.mp4
│       ├── parasite.mp4
│       └── dune.mp4
```

## Способы загрузки медиафайлов

### Способ 1: Через MinIO Web UI (самый простой)

1. Откройте MinIO Console: http://localhost:9090
2. Войдите с учетными данными:
   - Username: `minio`
   - Password: `minio123`
3. Выберите bucket `media`
4. Нажмите "Upload" и загрузите файлы:
   - Изображения в папку `frames/`
   - Видео в папку `videos/`

### Способ 2: Через Docker (Windows PowerShell)

```powershell
# Убедитесь, что контейнеры запущены
docker-compose up -d

# Запустите скрипт загрузки
.\scripts\upload-media.ps1
```

### Способ 3: Через Docker (Linux/Mac)

```bash
# Убедитесь, что контейнеры запущены
docker-compose up -d

# Сделайте скрипт исполняемым
chmod +x scripts/upload-media-docker.sh

# Запустите скрипт
./scripts/upload-media-docker.sh
```

### Способ 4: Через MinIO Client (mc)

Если у вас установлен MinIO Client:

```bash
# Настройка
mc alias set localminio http://localhost:9000 minio minio123

# Создание bucket
mc mb localminio/media

# Загрузка файлов
mc cp media/frames/* localminio/media/frames/
mc cp media/videos/* localminio/media/videos/
```

## Требования к файлам

### Изображения (frames)
- Формат: JPG, PNG, WebP
- Рекомендуемый размер: до 2MB
- Разрешение: 1920x1080 или меньше

### Видео (videos)
- Формат: MP4 (H.264)
- Рекомендуемый размер: до 50MB
- Длительность: 10-30 секунд
- Разрешение: 1280x720 или меньше

## Проверка загрузки

После загрузки проверьте:

1. Через MinIO Console: http://localhost:9090
2. Через API: http://localhost:8080/api/media/media/frames/matrix.jpg
3. В логах backend не должно быть ошибок при запросе медиа

## Примеры вопросов в БД

Убедитесь, что пути в БД соответствуют загруженным файлам:

```sql
-- Изображения
INSERT INTO questions (type, image_url, options, correct_answer) VALUES
('frame', 'frames/matrix.jpg', ARRAY['Матрица', 'Терминатор'], 'Матрица');

-- Видео
INSERT INTO questions (type, video_url, options, correct_answer) VALUES
('video', 'videos/avengers.mp4', ARRAY['Мстители', 'Тор'], 'Мстители');
```

## Решение проблем

### Файлы не отображаются

1. Проверьте, что файлы загружены в MinIO
2. Проверьте пути в БД (должны быть без начального `/`)
3. Проверьте консоль браузера на ошибки 404
4. Убедитесь, что backend может подключиться к MinIO

### Кнопка "Ответить" не работает

1. Откройте консоль браузера (F12)
2. Проверьте ошибки JavaScript
3. Убедитесь, что вопрос загружен (проверьте `question.value` в консоли)
4. Проверьте, что `sessionId` установлен

### Ошибка CORS

Если видите ошибки CORS, убедитесь, что:
- Backend запущен и доступен
- CORS middleware настроен правильно
- Frontend делает запросы на правильный URL


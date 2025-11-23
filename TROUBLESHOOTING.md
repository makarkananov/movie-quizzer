# Решение проблем

## Кнопки не кликаются

### Проблема
Кнопки вариантов ответа не реагируют на клики.

### Решения

1. **Откройте консоль браузера (F12)** и проверьте:
   - Есть ли ошибки JavaScript
   - Появляются ли логи при клике ("Клик по варианту:")
   - Проверьте значения: `answered`, `loading`, `sessionId`, `question`

2. **Проверьте состояние компонента:**
   ```javascript
   // В консоли браузера
   console.log('answered:', answered.value)
   console.log('loading:', loading.value)
   console.log('sessionId:', sessionId.value)
   console.log('question:', question.value)
   ```

3. **Перезагрузите страницу** и попробуйте снова

4. **Проверьте, что вопрос загружен:**
   - Должны быть видны варианты ответов
   - Таймер должен работать

## Файлы не отображаются

### Проблема
Изображения и видео не загружаются, показывается ошибка 404.

### Решения

1. **Проверьте загрузку файлов в MinIO:**
   ```bash
   # Проверьте через скрипт
   .\scripts\upload-media.ps1
   
   # Или через MinIO Console
   # Откройте http://localhost:9090
   # Проверьте bucket "media"
   # Файлы должны быть в папках:
   #   - media/frames/matrix.jpg
   #   - media/videos/avengers.mp4
   ```

2. **Проверьте пути в БД:**
   ```sql
   SELECT image_url, video_url FROM questions LIMIT 5;
   ```
   Пути должны быть: `frames/matrix.jpg` (без префикса `media/`)

3. **Проверьте логи backend:**
   ```bash
   docker-compose logs backend | grep Media
   ```
   Должны быть логи вида: `Media request: bucket=media, file=frames/matrix.jpg`

4. **Проверьте структуру в MinIO:**
   - Bucket: `media`
   - Путь к файлу: `frames/matrix.jpg` (не `media/frames/matrix.jpg`)
   - Путь к видео: `videos/avengers.mp4` (не `media/videos/avengers.mp4`)

### Правильная структура в MinIO

```
bucket: media
├── frames/
│   ├── matrix.jpg
│   ├── inception.jpg
│   ├── pulp_fiction.jpg
│   ├── shawshank.jpg
│   └── godfather.jpg
└── videos/
    ├── avengers.mp4
    ├── interstellar.mp4
    ├── joker.mp4
    ├── parasite.mp4
    └── dune.mp4
```

## Вопросы с цитатами показывают кнопки вместо поля ввода

Это **правильное поведение**! В seed.sql все вопросы (включая цитаты) имеют варианты ответов (`options`).

Если вы хотите, чтобы цитаты использовали поле ввода, нужно:
1. Убрать `options` из вопросов типа `quote` в seed.sql
2. Или изменить логику во frontend

## Проверка работы

### 1. Проверка БД
```bash
docker exec -it quiz-postgres psql -U app -d quiz -c "SELECT type, image_url, video_url, array_length(options, 1) as options_count FROM questions LIMIT 5;"
```

### 2. Проверка MinIO
```bash
# Через Docker
docker run --rm --network movie-quizzer_default minio/mc ls localminio/media/frames/
docker run --rm --network movie-quizzer_default minio/mc ls localminio/media/videos/
```

### 3. Проверка API
```bash
# Проверьте доступность медиа
curl http://localhost:8080/api/media/media/frames/matrix.jpg
```

### 4. Проверка frontend
- Откройте консоль браузера (F12)
- Перейдите на вкладку Network
- Попробуйте загрузить вопрос
- Проверьте запросы к `/api/game/sessions` и `/api/media/...`

## Частые ошибки

### Ошибка: "file is empty or not found"
- Файл не загружен в MinIO
- Неправильный путь к файлу
- Файл действительно пустой

### Ошибка: "options is not defined"
- Проблема с сериализацией JSON из backend
- Проверьте, что backend пересобран после изменений

### Ошибка: "Cannot read property 'id' of null"
- Вопрос не загружен
- Проверьте запрос к `/api/game/sessions`


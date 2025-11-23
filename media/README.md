# Медиафайлы для игры

Эта папка содержит медиафайлы (изображения и видео), которые нужно загрузить в MinIO.

## Структура

```
media/
├── frames/          # Изображения кадров из фильмов
│   ├── matrix.jpg
│   ├── inception.jpg
│   ├── pulp_fiction.jpg
│   ├── shawshank.jpg
│   └── godfather.jpg
└── videos/          # Видеофрагменты из фильмов
    ├── avengers.mp4
    ├── interstellar.mp4
    ├── joker.mp4
    ├── parasite.mp4
    └── dune.mp4
```

## Как добавить файлы

1. Поместите изображения в папку `frames/`
2. Поместите видео в папку `videos/`
3. Убедитесь, что имена файлов совпадают с путями в `backend/scripts/seed.sql`

## Загрузка в MinIO

### Автоматически (через скрипт):

```bash
# Windows PowerShell
.\scripts\upload-media.ps1

# Linux/Mac
./scripts/upload-media-docker.sh
```

### Вручную (через MinIO Console):

1. Откройте http://localhost:9090
2. Войдите: `minio` / `minio123`
3. Выберите bucket `media`
4. Загрузите файлы:
   - Изображения в папку `frames/`
   - Видео в папку `videos/`

## Требования к файлам

- **Изображения**: JPG, PNG (рекомендуется до 2MB)
- **Видео**: MP4 (рекомендуется до 50MB, длительность 10-30 сек)

## Важно

- Имена файлов должны точно совпадать с путями в БД (seed.sql)
- Файлы должны быть загружены в bucket `media` в MinIO
- Пути в БД: `frames/matrix.jpg` (без префикса `media/`)


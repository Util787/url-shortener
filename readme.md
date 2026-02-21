## Краткая документация API

Небольшая инструкция по запуску и использованию HTTP API сервиса сокращения ссылок.

### Коротко
- Базовый префикс: `/`
- Формат: JSON для POST-запросов

### Переменные окружения
Скопируйте `.env.example` в `.env` и заполните значения перед запуском.

Основные переменные (в файле `.env.example`):

- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` — доступ к PostgreSQL
- `HTTP_SERVER_HOST`, `HTTP_SERVER_PORT` — адрес и порт HTTP сервера
- `TG_BOT_TOKEN` — токен телеграм-бота (опционально)

### Запуск
1. Запустите PostgreSQL (например, в Docker) и примените миграции из `migrations/postgres`.
# Эндпоинты

Короткое описание HTTP‑эндпоинтов сервиса (маршруты находятся в корне сервера).

- POST /save
	- Body (JSON): { "long_url": "https://example.com/..." }
	- Ответ 200: { "message": "URL saved successfully with short URL: http://<host>:<port>/<id>" }

- GET /random
	- Ответ 200: { "short_url": "http://<host>:<port>/<id>" }

- GET /:short_url_id
	- Пример: GET /5vrTql7AOt
	- Возвращает HTTP 302 редирект на оригинальную длинную ссылку

- POST /delete
	- Body (JSON): { "short_url": "http://<host>:<port>/<id>" } (или используйте `/delete/:short_url` если предпочитаете path param)
	- Ответ 200: { "message": "URL deleted successfully" }

---

Файл с примерами переменных окружения: `.env.example`.



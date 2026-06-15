# GeoLogicApi

Backend-сервис для хранения пользователей, локаций, событий и вычисления рейтинга местности при выборе места для открытия бизнеса.

## Кратко о проекте

GeoLogicApi предоставляет REST API для:
- хранения пользователей и их связей с локациями;
- хранения локаций и инфраструктуры (с геокоординатами);
- регистрации событий, привязанных к точкам на карте;
- расчёта и выдачи рейтинга местности для оценки привлекательности участка под бизнес.

Архитектура ориентирована на DDD: `handler` → `service` → `repository` с DI

## Стек технологий

- Go (1.25)
- Gin (HTTP)
- GORM (ORM)
- PostgreSQL + PostGIS (геоданные)
- log/slog (логирование)
- github.com/caarlos0/env + godotenv (конфигурация)
- testcontainers-go (e2e)
- stretchr/testify (assertions)

## Схема данных (план)

### Users
- `id` — PK
- `tg_id` — int

### Events
- `id` — PK
- `latitude` — geography
- `longitude` — geography
- `date` — datetime
- `info` — string

### Locations
- `id` — PK
- `latitude` — geography
- `longitude` — geography
- `address` — string
- `type` — string
- `info` — string

### User_locations
- `id` — PK
- `user_id` — int (FK → Users.id)
- `location_id` — int (FK → Locations.id)

### Infrastructure
- `id` — PK
- `latitude` — geography
- `longitude` — geography
- `address` — string
- `type` — string
- `info` — string

## Что уже реализовано

- Базовая структура проекта с разделением слоёв (`internal/api`, `internal/handler`, `internal/service`, `internal/repository`, `internal/database`).
- CRUD для пользователей (handler/service/repository).
- Подключение к PostgreSQL + PostGIS, автоматическая миграция модели пользователя.
- Централизованная модель ошибок (`internal/errors` → `AppError`) и middleware для ответа клиенту.
- Логирование запросов через `internal/logger`.
- E2E-тесты с `testcontainers-go` (контейнер Postgres).

## Что планируется сделать

- Добавить модели и репозитории для `Locations`, `Events`, `Infrastructure`, `User_locations`.
- Реализовать геопространственные запросы с использованием PostGIS
- Реализовать алгоритм расчёта рейтинга местности (агрегация близости инфраструктуры, событий, плотности и т.п.).
- Добавить аутентификацию/авторизацию (если потребуется).
- Написать unit-тесты для сервиса и репозиториев, расширить интеграционные тесты.
- Сделать API документацию (OpenAPI/Swagger).



## Конфигурация и запуск

1. Скопировать шаблон окружения:
```bash
cp .env.template .env
```

2. Настроить переменные в .env (DB, SERVER_PORT, LOG_LEVEL и т.п.).

3. Запуск локально через докер и task:
```bash
task db:up
task run
```
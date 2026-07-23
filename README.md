# GeoLogicApi

Backend-сервис для хранения пользователей, локаций, событий, аренды недвижимости и вычисления рейтинга местности при выборе места для открытия бизнеса.

## Кратко о проекте

GeoLogicApi предоставляет REST API для:
- управления пользователями и их связями с локациями;
- хранения объектов коммерческой аренды с геопривязкой;
- хранения инфраструктуры города и ключевых точек притяжения;
- регистрации событий, привязанных к точкам на карте;
- расчёта и выдачи рейтинга местности для оценки привлекательности участка под бизнес.

Архитектура проекта строго ориентирована на принципы Clean Architecture / DDD: `handler` (транспорт) → `service` (бизнес-логика) → `repository` (данные) с управлением зависимостями через DI.

## Стек технологий

- Go (1.26.1)
- Gin (HTTP-роутинг и валидация)
- GORM (ORM)
- PostgreSQL + PostGIS (хранение и индексация геоданных)
- log/slog (структурированное логирование)
- github.com/caarlos0/env + godotenv (конфигурация)
- testcontainers-go (e2e-тестирование в изолированных контейнерах)
- stretchr/testify (пакеты assert и require для тестирования)

## Схема данных (план)

### Users
- `id` — PK (uint64)
- `tg_id` — int
- `created_at` — datetime
- `updated_at` — datetime

### Rents (Коммерческая аренда)
- `id` — PK (uint64)
- `location` — geometry(Point, 4326) [not null]
- `address` — string (varchar 255)
- `info` — string (varchar 255, nullable)
- `created_at` — datetime
- `updated_at` — datetime

### Events
- `id` — PK
- `location` — geometry(Point, 4326)
- `date` — datetime
- `info` — string
- `created_at` — datetime
- `updated_at` — datetime

### TrakedLocations
- `id` — PK
- `user_id` — int (FK → Users.id)
- `location` — geometry(Point, 4326)
- `created_at` — datetime
- `updated_at` — datetime

### Infra
- `id` — PK
- `location` — geometry(Point, 4326)
- `address` — string
- `type` — string
- `info` — string
- `created_at` — datetime
- `updated_at` — datetime

### InfraType
- `id` — PK
- `slug` — string, unique identifier for code usage (e.g. `subway`, `cafe`)
- `name` — string, display name (e.g. `Метро`)
- `weight` — float, base score coefficient
- `max_radius` — uint16, maximum influence radius in meters
- `created_at` — datetime
- `updated_at` — datetime

## Что уже реализовано

- **Базовая структура**: Архитектура проекта с разделением слоёв (`internal/api`, `internal/handler`, `internal/service`, `internal/repository`, `internal/database`).
- **Пользователи**: Полный CRUD для пользователей (handler/service/repository) с автоматической миграцией.
- **Коммерческая аренда (Rents)**: 
  - Разработан безопасный PATCH-апдейт на указателях для частичного обновления данных.
  - Реализован геопространственный поиск доступных объектов в радиусе с помощью функций PostGIS (`ST_DWithin` с кастингом в `geography`).
  - Реализована строгая валидация входящих Query-параметров (координаты, ограничения радиуса) на уровне Gin Binding.
- **Инфраструктура**: Подключение к PostgreSQL + PostGIS, централизованная модель ошибок (`internal/errors` → `AppError`) и middleware для ответов клиенту.
- **Логирование**: Логирование запросов и SQL-операций через встроенный `internal/logger` (`slog`).
- **Тестирование**: Мощное покрытие E2E-тестами (`testcontainers-go`). Написаны табличные тесты (Table-Driven Tests) для верификации успешных PATCH-обновлений, гео-поиска `Available`, а также негативные тесты на валидацию «грязных» данных и сломанного JSON.

## Что планируется сделать

- Добавить модели и репозитории для оставшихся сущностей: `Events`, `Infra`, `Users_locations`.
- Создать пространственные индексы `GIST` для высокопроизводительных гео-запросов.
- Реализовать алгоритм расчёта рейтинга местности (агрегация близости инфраструктуры, событий, плотности и т.п.).
- Добавить аутентификацию/авторизацию (если потребуется).
- Расширить unit-тестирование бизнес-логики сервисов.
- Подготовить API документацию (OpenAPI/Swagger).
- Расширить документацию по новым сущностям: `Events`, `Infra`, `Users_locations`, `InfraType`.

## Конфигурация и запуск

1. Скопировать шаблон окружения:
```bash
cp .env.template .env
```

2. Настроить переменные в `.env` (DB, SERVER_PORT, LOG_LEVEL и т.п.).

3. Запуск локально базы данных и приложения через `task`:
```bash
task db:up
task run
```
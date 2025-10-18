# Auth Microservice

Микросервис аутентификации с поддержкой JWT токенов (access + refresh).

## 📖 Описание

Микросервис для регистрации, аутентификации и тд. Реализован с использованием чистой архитектуры (Clean Architecture) и современных практик разработки на Go.

## 🛠 Технологии

- **Go 1.23** - язык программирования
- **gRPC** - для межсервисной коммуникации
- **PostgreSQL 14** - база данных
- **JWT (HS256)** - токены аутентификации
- **Bcrypt** - хеширование паролей
- **Goose** - миграции БД
- **Uber FX** - dependency injection
- **pgx/v4** - драйвер PostgreSQL с connection pooling
- **Uber Ratelimit - Либа для ограничения запросов
- **Viper - Либа для конфигурации, гибчее и чище чем стандартная 

**Принципы:**
- Clean Architecture
- Dependency Injection (Uber FX)
- Repository Pattern
- Разделение ответственности

## 🚀 Быстрый старт

### Требования

- Go 1.23+
- Docker & Docker Compose
- Make

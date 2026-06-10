# TundraMarket-Api

# Настройка переменных окружения

Перед запуском приложения необходимо создать файл `.env` в корне проекта.

## 1. Создайте файл `.env`

Скопируйте содержимое файла `.env.example`:

```bash
cp .env.example .env
```

Или создайте файл вручную.

## 2. Заполните переменные окружения

Пример содержимого файла `.env`:

```env
# Postgres
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DATABASE=tundra_market
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin

DATABASE_URL=postgres://admin:admin@postgres:5432/tundra_market?sslmode=disable

AUTH_TOKEN_SECRET=secret
AUTH_TOKEN_TTL=8760h
```

## Описание переменных

| Переменная          | Описание                                                                                   |
| ------------------- | ------------------------------------------------------------------------------------------ |
| `POSTGRES_HOST`     | Хост PostgreSQL. При использовании Docker Compose обычно `postgres`.                       |
| `POSTGRES_PORT`     | Порт PostgreSQL. По умолчанию `5432`.                                                      |
| `POSTGRES_DATABASE` | Название базы данных.                                                                      |
| `POSTGRES_USER`     | Пользователь базы данных.                                                                  |
| `POSTGRES_PASSWORD` | Пароль пользователя базы данных.                                                           |
| `DATABASE_URL`      | Строка подключения к PostgreSQL.                                                           |
| `AUTH_TOKEN_SECRET` | Секретный ключ для подписи JWT-токенов. В продакшене используйте случайную длинную строку. |
| `AUTH_TOKEN_TTL`    | Время жизни токена. Например: `24h`, `168h`, `8760h`.                                      |

## Рекомендации для продакшена

* Используйте сложный пароль для `POSTGRES_PASSWORD`.
* Замените `AUTH_TOKEN_SECRET` на случайный секрет длиной не менее 32 символов.

## Запуск проекта через Docker

1. Соберите и запустите контейнеры

В корне проекта выполните:
docker compose up -d --build

Параметры команды:
-d — запуск в фоновом режиме.
--build — принудительная пересборка образов.

3. Проверьте состояние контейнеров
docker compose ps

Для просмотра логов:
docker compose logs -f

4. Остановка проекта

Остановить контейнеры:
docker compose stop

Остановить и удалить контейнеры:
docker compose down

Остановить и удалить контейнеры вместе с томами:
docker compose down -v



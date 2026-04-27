# Todo API

REST API на Go с PostgreSQL и Docker.

## Стек

- **Go** — язык программирования
- **PostgreSQL** — база данных
- **Docker + docker-compose** — контейнеризация
- **gorilla/mux** — HTTP роутер

## Запуск

```bash
docker-compose up --build
```

API будет доступен на `http://localhost:8080`

## Эндпоинты

| Метод | URL | Описание |
|-------|-----|----------|
| GET | /todos | Получить все задачи |
| POST | /todos | Создать задачу |
| PUT | /todos/{id} | Обновить задачу |
| DELETE | /todos/{id} | Удалить задачу |

## Примеры запросов

### Создать задачу
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Изучить Go","completed":false}'
```

### Получить все задачи
```bash
curl http://localhost:8080/todos
```

### Обновить задачу
```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Изучить Go","completed":true}'
```

### Удалить задачу
```bash
curl -X DELETE http://localhost:8080/todos/1
```

## Структура проекта

```
todo-api/
├── main.go          # Точка входа, роутер
├── handler/
│   └── todo.go      # HTTP обработчики
├── model/
│   └── todo.go      # Структура Todo
├── db/
│   └── postgres.go  # Подключение к PostgreSQL
├── Dockerfile       # Multi-stage сборка
├── docker-compose.yml
└── go.mod
```
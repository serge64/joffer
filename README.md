# Joffer (пет-проект)

Сервис для массового поиска вакансий с различных рекрутмент-площадок и дальнейшего отклик на них.
На текущий момент подключен hh.ru.

## Требования

1. Go (1.16 и выше)
2. [Golang-migrate](https://github.com/golang-migrate/migrate)
2. Docker
3. Docker-compose

## Установка и запуск

```bash
make compose_up
make migrate_up
make build
bin/joffer
```

## TODO

1. Добавить тесты.
2. Добавить еще больше тестов.
3. Подключить другие площадки.
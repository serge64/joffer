# Joffer (пет-проект)

Приложение для массового поиска вакансий с различных рекрутмент-площадок и дальнейшего отклика на них.
На текущий момент подключен hh.ru.

## Требования

1. Go (1.16 и выше)
2. [Golang-migrate](https://github.com/golang-migrate/migrate)
2. Docker
3. Docker-compose

## Сборка

```bash
$ make compose_up
$ make migrate_up
$ make build
```

## Запуск

```bash
$ bin/joffer
```

Приложение будет запущено на порту 8080.

## TODO

1. Добавить тесты.
2. Добавить еще больше тестов.
3. Подключить другие площадки.
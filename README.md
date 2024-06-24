# TalkBoard
a system for creating and reading posts and comments.

## Содержание
- [TalkBoard](#TalkBoard)
- [содержание](#содержание)
- [Архитектура](#архитектура)
- [Локальный запуск](#локальный-запуск)

## Архитектура

Проект основан на "чистой" архитектуре с использованием GraphQL и разработан на языке Go. Реализованная архитектура включает три основных слоя: handler, service и repository. В слое handler обрабатываются GraphQL запросы, слой service содержит бизнес-логику, а в repository выполняются операции с базой данных или с памятью (это зависит от переданного параметра при запуске программы).

Для упрощения развертывания и изоляции окружений используется Docker с docker-compose для сборки и запуска контейнеров. Это позволяет легко воспроизводить среду разработки и обеспечивать консистентность окружений между разработкой и продуктивной средой.

Планируется дальнейшее развитие проекта с акцентом на улучшение производительности и надежности системы.

## Локальный запуск
Для развертывания сервиса необходимы [Docker](https://docs.docker.com/engine/install/)

**Важно:** Перед запуском убедитесь, что вы установили пароли. Для этого создайте файл .env в корневой директории и добавьте туда:\
`DB_PASSWORD=...`

1. `git clone --recurse-submodules https://github.com/klausfun/TalkBoard`
2. `cd TalkBoard`
3. `docker pull postgres` - для скачивания образа postgres
4. `$env:STORAGE_TYPE = ...` вместо ... нужно указать `"memory"` или `"postgres"`, в ином случае будет использоваться postgres (по умолчанию)
5. `docker run --name=talkBoard-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres` (вместо 'qwerty' необходимо указать Ваш пароль от БД)
6. `migrate -path ./schema_db -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up` (вместо 'qwerty' необходимо указать Ваш пароль от БД)
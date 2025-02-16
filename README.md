# запуск сервера
```
docker-compose up --build
```
# миграции
```
docker-compose exec app goose -dir migrations postgres "user=avito password=secret dbname=avito_shop host=postgres sslmode=disable" up
```
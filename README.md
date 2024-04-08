# Cars Catalog API
Effective mobile test task. Cars catalog API

# Start
With Docker: \
```docker compose up``` \
Manually: \
```go build effective_mobile_test/cmd/cars-catalog``` \
Set env CONFIG_PATH={Yours GOPATH}/effective_mobile_test/.env \
```./cars-catalog```

# Migrations
Using golang-migrate \
```migrate -path ./internal/storage/schema -database "postgres://root:root@localhost:5432/cars_catalog?sslmode=disable" up```

# Swagger
```http://localhost:8082/swagger/index.html#```

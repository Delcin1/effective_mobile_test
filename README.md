# Cars Catalog API
Effective mobile test task. Cars catalog API

# Start
**Help API URL must be changed in .env** \
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

# ТЗ
Реализовать каталог автомобилей. Необходимо реализовать следующее
1. Выставить rest методы
	1. Получение данных с фильтрацией по всем полям и пагинацией 
	2. Удаления по идентификатору
	3. Изменение одного или нескольких полей по идентификатору
	4. Добавления новых автомобилей в формате
```json
{
    "regNums": ["X123XX150"] // массив гос. номеров
}
```
2. При добавлении сделать запрос во внешнее АПИ, описанного сваггером (это описание некоторого внешнего АПИ, которого нет, но к которому надо обращаться. Реализованное, согласно описанию, АПИ будет использоваться при проверке)

```yaml
openapi: 3.0.3
info:
  title: Car info
  version: 0.0.1
paths:
  /info:
    get:
      parameters:
        - name: regNum
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Ok
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Car'
        '400':
          description: Bad request
        '500':
          description: Internal server error
components:
  schemas:
    Car:
      required:
        - regNum
        - mark
        - model
        - owner
      type: object
      properties:
        regNum:
          type: string
          example: X123XX150
        mark:
          type: string
          example: Lada
        model:
          type: string
          example: Vesta
        year:
          type: integer
          example: 2002
        owner:
          $ref: '#/components/schemas/People'
    People:
      required:
        - name
        - surname
      type: object
      properties:
        name:
          type: string
        surname:
          type: string
        patronymic:
          type: string
```
3. Обогащенную информацию положить в БД postgres (структура БД должна быть создана путем миграций при старте сервиса)
4. Покрыть код debug- и info-логами
5. Вынести конфигурационные данные в .env-файл
6. Сгенерировать сваггер на реализованное АПИ


# Readme к реализованным запросам

## API Endpoints

| Endpoint                    | Описание                                  |
|-----------------------------|-------------------------------------------|
| `POST /pickup_point`        | Добавляет новый ПВЗ                       |
| `PUT /pickup_point`         | Изменяет существующий ПВЗ                 |
| `GET /pickup_point/[id]`    | Возвращает информацию о ПВЗ с заданным id |
| `DELETE /pickup_point/[id]` | Удаляет ПВЗ с заданным id                 |
| `GET /pickup_point/list`    | Возвращает список всех ПВЗ                |

## Примеры использования

### Добавление ПВЗ

```bash
 curl https://localhost:9000/pickup_point -k --cacert ./server.crt -u test:test -i -d \
 '{
    "name": "PickupPoint_1",
    "address":"Address_1",
    "phone_number":"+7-999-999-99-99"
  }'
```
### Обновление информации о ПВЗ

```bash
curl https://localhost:9000/pickup_point -k --cacert ./server.crt -X PUT -u test:test -i -d \
'{
    "id":1,
    "name": "Updated_PickupPoint_1",
    "address":"Updated_Address_1",
    "phone_number":"+7-999-999-99-99"
}'
```

### Получение информации о ПВЗ

```bash
curl https://localhost:9000/pickup_point/1 -k --cacert ./server.crt -u test:test -i
```

### Удаление ПВЗ

```bash
curl -X DELETE https://localhost:9000/pickup_point/1 -k --cacert ./server.crt -u test:test -i
```

### Получение списка всех ПВЗ
```bash
curl https://localhost:9000/pickup_point/list -k --cacert ./server.crt -u test:test -i 
```

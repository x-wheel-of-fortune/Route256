# Readme к реализованным запросам

## API Endpoints

| Endpoint            | Описание                                  |
|---------------------|-------------------------------------------|
| `AddPickupPoint`    | Добавляет новый ПВЗ                       |
| `UpdatePickupPoint` | Изменяет существующий ПВЗ                 |
| `GetPickupPoint`    | Возвращает информацию о ПВЗ с заданным id |
| `DeletePickupPoint` | Удаляет ПВЗ с заданным id                 |
| `ListPickupPoint`   | Возвращает список всех ПВЗ                |

## Примеры использования

### Добавление ПВЗ

```bash
AddPickupPoint
'{
    "pickupPoint": {
        "name":"Name",
        "address":"Address",
        "phone_number": "PhoneNumber"
    }
}'
```

### Обновление информации о ПВЗ

```bash
UpdatePickupPoint
'{
    "pickupPoint": {
        "id":1,
        "name":"Name1",
        "address":"Address1",
        "phone_number": "PhoneNumber1"
    }
}'
```

### Получение информации о ПВЗ

```bash
GetPickupPoint
{
    "id":1
}
```

### Удаление ПВЗ

```bash
DeletePickupPoint
{
    "id":1
}
```

### Получение списка всех ПВЗ

```bash
ListPickupPoint
{}
```

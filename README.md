
# Key Value Store

This is a non-persisitent key value store, that uses sqlite as it's database

## Installation

This project is built with Golang and requires sqlite3

```bash
  $ make
```

## Usage/Examples

There are three api methods \
Add a value
POST http://{API_URL}/v1/key
Payload:
```json
{
    "key": "abc",
    "value":"xyz",
    "ttl": 3600 // value in milliseconds
}
```

Get a value
POST http://{API_URL}/v1/key
Payload:
```json
{
    "key": "abc",
}
```

Delete a value
PATCH http://{API_URL}/v1/key
Payload:
```json
{
    "key": "abc",
}
```
## Archtecture

![Architecture](https://github.com/anmolgaud/kv-store-golang/blob/main/images/arch.png?raw=true)


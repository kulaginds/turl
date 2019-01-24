# Tiny URL
Service for short url.

## Installation
1. `git clone https://github.com/kulaginds/turl.git && cd turl`
1. Create MySQL database with user.
1. Copy `.env-example` to `.env` and edit DB credentials.

## Build
```bash
go build -o tiny_url
```

## Run
```bash
export $(cat .env | xargs) && ./tiny_url
```

## Test queries
1. Short:
```bash
curl -i -X POST -d '{"url":"http://google.ru"}' http://localhost:8080/short
```
1. Long:
```bash
curl -i -X POST -d '{"url":"http://localhost:8080/baaa"}' http://localhost:8080/long
```
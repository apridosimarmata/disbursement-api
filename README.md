# Disbursement API
#### submitted by Imam Aprido Simarmata

## Depedencies

### 1. Postgres

Used to store disbursements data

### 2. Redis

Used to store tokens as key and wallet id as value\
Also provide distributed lock (in case locks are being used in multiple pod/machine) to prevent race condition on `deposits` and `withdrawal`\
Redlock (implemented with redsync) also provide TTL for each lock to prevent deadlock.\

## How to setup


## How to setup

1. Clone this repository
2. Make sure you have docker installed on your machine
3. Get in to the directory `cd disbursement`
4. Run `docker-compose up`
5. Open new terminal and run `docker exec -it disbursement sh`
6. Once you are already in the container shell, run this command:

`cd /go/src/disbursement && go install github.com/pressly/goose/v3/cmd/goose@v3.15.0 && export PATH="$PATH:$HOME/go/bin"&& goose -dir infrastructure/migrations postgres "host=postgres port=5432 user=postgres password=postgres dbname=disbursement sslmode=disable" up`

## Service Description

### Account

Path -> `/api/v1/accounts\

GET `/api/v1/accounts/{number}`\


```
GET /api/v1/accounts/211833558

{
    "status": "success",
    "data": {
        "name": "name 1",
        "number": "211833558"
    }
}
```


```
GET /api/v1/accounts/211833553

{
    "status": "fail",
    "data": {
        "error": "account not found"
    }
}
```

## Happy testing :)
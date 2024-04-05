# Disbursement API
#### submitted by Imam Aprido Simarmata

## Depedencies

### 1. Postgres

Used to store disbursements data


## How to setup

1. Clone this repository
2. Make sure you have docker installed on your machine
3. Get in to the directory `cd disbursement`
4. Run `docker-compose up`
5. Open new terminal and run `docker exec -it disbursement sh`
6. Once you are already in the container shell, run this command:

`cd /go/src/disbursement && go install github.com/pressly/goose/v3/cmd/goose@v3.15.0 && export PATH="$PATH:$HOME/go/bin"&& goose -dir infrastructure/migrations postgres "host=postgres port=5432 user=postgres password=postgres dbname=disbursement sslmode=disable" up`

## Service Description

### Mock server

The mock server was created with Postman Mockserver

url: https://91d6ea63-ae7f-4167-936f-699986c9fa36.mock.pstmn.io/api/v1/

documention: https://documenter.getpostman.com/view/34064887/2sA35LVzKC

### Disbursement API

documentation: https://documenter.getpostman.com/view/34064887/2sA35LVzKA


## Happy testing :)
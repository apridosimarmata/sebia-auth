# Mini Wallet
#### submitted by Imam Aprido Simarmata

## Depedencies

### 1. Postgres

Used to store wallet & transaction data

### 2. Redis

Used to store tokens as key and wallet id as value\
Also provide distributed lock (in case locks are being used in multiple pod/machine) to prevent race condition on `deposits` and `withdrawal`\
Redlock (implemented with redsync) also provide TTL for each lock to prevent deadlock.\

## How to setup


## How to setup

1. Clone this repository
2. Make sure you have docker installed on your machine
3. Get in to the directory `cd mini-wallet`
4. Run `docker-compose up`
5. Open new terminal and run `docker exec -it mini-wallet sh`
6. Once you are already in the container shell, run this command:

`cd /go/src/mini-wallet && go install github.com/pressly/goose/v3/cmd/goose@v3.15.0 && export PATH="$PATH:$HOME/go/bin"&& goose -dir infrastructure/migrations postgres "host=postgres port=5432 user=postgres password=postgres dbname=mini-wallet sslmode=disable" up`

## Happy testing :)
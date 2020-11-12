# Rate-Limiter
A simple homemade rate limiter
## Installation and Building

```bash
# clone repository
git clone https://github.com/johnhckuo/Rate-Limiter.git
cd Rate-Limiter

# building docker image
make docker

# running redis & golang application locally
docker-compose up

```

## Environment variables

RATE_LIMITER
> It is used to decide which rate limiter to use, default value is `REDIS_LIMITER`

PERSIST_STORAGE
> It's used to decide which storage type to use, default value is `REDIS_STORAGE`

DB_CONNECTION_STRING=redis://localhost:6379
> DB connection string

RATE_LIMIT  
> This is the number of request we can allow per `RATE_LIMIT_EXPIRATION_SECOND`, default value is 60

BURST_LIMIT
> The number of burst request within 1 second, default value is 10

RATE_LIMIT_EXPIRATION_SECOND
> Time interval of `RATE_LIMIT`, default value is 60 (seconds)


## Commands

```bash
# running the application locally
make run

# start sending request to test the behavior of rate limiter
make test-api

# build image
make docker

```

## Powered by
- [Redis](https://redis.io/)

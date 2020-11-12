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

- RATE_LIMITER=REDIS_LIMITER
- PERSIST_STORAGE=REDIS_STORAGE
- DB_CONNECTION_STRING=redis://localhost:6379
- RATE_LIMIT=60
- RATE_LIMIT_EXPIRATION_SECOND=60
- BURST_LIMIT=10

## Commands

```bash
# running the application locally
make run

# starting sending request to test the behavior of rate limiter
make test-api

# build image
make docker

```

## Powered by
- [Redis](https://redis.io/)

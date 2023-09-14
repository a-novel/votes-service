# Votes service

Manage voting system.

## Prerequisites

- Download [Go](https://go.dev/doc/install)
- Install [Mockery](https://vektra.github.io/mockery/latest/installation/)
- Clone [go-framework](https://github.com/a-novel/go-framework)
    - From the framework, run `docker compose up -d`

## Installation

Create a env file.

```bash
touch .envrc
```
```bash
printf 'export POSTGRES_URL="postgres://votes@localhost:5432/agora_votes?sslmode=disable"
export POSTGRES_URL_TEST="postgres://test@localhost:5432/agora_votes_test?sslmode=disable"
' > .envrc
```
```bash
direnv allow .
```

Set the database up.
```bash
make db-setup
```

## Commands

### Run the API

```bash
make run
```
```bash
curl http://localhost:2042/ping
# Or curl http://localhost:2042/healthcheck
```

### Run tests

```bash
make test
```

### Update mocks

```bash
mockery
```

### Open a postgres console

```bash
make db
# Or make db-test
```

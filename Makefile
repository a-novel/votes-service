COVER_FILE=$(CURDIR)/coverage.out
BIN_DIR=$(CURDIR)/bin

PKG="github.com/a-novel/votes-service"

PKG_LIST=$(shell go list $(PKG)/... | grep -v /vendor/)

# Runs the test suite.
test:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --junitfile report.xml --format pkgname -- -count=1 -p 1 -v -coverpkg=./...

# Runs the test suite in race mode.
race:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --format pkgname -- -race -count=1 -p 1 -v -coverpkg=./...

# Run the test suite in memory-sanitizing mode. This mode only works on some Linux instances, so it is only suitable
# for CI environment.
msan:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		env CC=clang env CXX=clang++ gotestsum --packages="./..." --format testname -- -msan -short $(PKG_LIST) -p 1

db-setup:
	psql -h localhost -p 5432 -U postgres agora -a -f init.sql

# Plugs into the development database.
db:
	psql -h localhost -p 5432 -U users agora_forum

# Plugs into the test database.
db-test:
	psql -h localhost -p 5432 -U test agora_forum_test

run:
	direnv allow . && source .envrc && go run ./cmd/api/main.go

.PHONY: all test race msan db db-test

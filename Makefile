# include variables from .env file
include ./.envrc

#========================================================================================================#
# HELPERS
#========================================================================================================#

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

#Create the new confirm target.
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


#=========================================================================================================#
# BUILD
#=========================================================================================================#

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
go_version = $(shell go version | awk -v N=3 '{print $$3}')
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description} -X main.goVersion=${go_version}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'building cmd/api...'
	go build -ldflags=${linker_flags} -o ./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o ./bin/linux_amd64/api ./cmd/api
	

#=========================================================================================================#
# DEVELOPMENT
#=========================================================================================================#

.PHONY: run/binary
run/binary:
	@./bin/api --db-dsn=${CRESTCOUNTRIES_DB_DSN} --db-max-open-conns=300 --db-max-idle-conns=300
# @./bin/api --db-dsn=${CRESTCOUNTRIES_DB_DSN}
	
## run/api: run the cmd/api application
.PHONY: run/api
launch/api:
	@go run cmd/api/* -db-dsn=${CRESTCOUNTRIES_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${CRESTCOUNTRIES_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${CRESTCOUNTRIES_DB_DSN} up


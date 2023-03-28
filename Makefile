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

ver?=0.0.6
IMAGE_NAME=ghcr.io/bxffour/crest/api
## build/docker: build the cmd/api dockerfile
.PHONY: build/docker
build/docker:
	docker build \
		--build-arg IMAGE_VERSION=${ver} \
		--build-arg IMAGE_REVISION=${git_description} \
		--build-arg IMAGE_CREATED=${current_time} \
		--build-arg LINKER_FLAGS=$(linker_flags) \
		--tag ${IMAGE_NAME}:${ver} \
		--tag ${IMAGE_NAME}:${ver}-${git_description} \
		.
		

#=========================================================================================================#
# DEVELOPMENT
#=========================================================================================================#

.PHONY: run/binary
run/binary:
	@./bin/api --config=./config.toml --secret=./secret.toml 
	
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


#=========================================================================================================#
# SSL
#=========================================================================================================#

CONFIG_PATH=${HOME}/.crest_test
WORKDIR=deployments/ssl

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: clean
clean:
	rm -rf ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert \
			-initca ${WORKDIR}/ca-csr.json | cfssljson -bare ca

	cfssl gencert \
			-ca=ca.pem \
			-ca-key=ca-key.pem \
			-config=${WORKDIR}/ca-config.json \
			-profile=server \
			${WORKDIR}/server-csr.json | cfssljson -bare pgsql

	cfssl gencert \
			-ca=ca.pem \
			-ca-key=ca-key.pem \
			-config=${WORKDIR}/ca-config.json \
			-profile=client \
			-cn="crest" \
			${WORKDIR}/client-csr.json | cfssljson -bare postgresql

	cp pgsql*.pem ca.pem deployments/postgres/bleh
	mv *.pem *.csr ${CONFIG_PATH}

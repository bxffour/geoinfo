# syntax=docker/dockerfile:1

## 
## BUILD
##
FROM golang:1.18.2-alpine3.15 as build-env

COPY . /go/src/app
WORKDIR /go/src/app

RUN go mod tidy
RUN go build -o /go/bin/crest-app ./cmd/api

##
## DEPLOY
##
FROM alpine:3.15.4
COPY --from=build-env /go/bin/crest-app /usr/local/bin/

ENV CREST_PORT=8080
ENV CREST_ENV=development
ENV CREST_DSN_PATH=/etc/crest/.env
ENV CREST_DB_MAX_OPEN_CONNS=25
ENV CREST_DB_MAX_IDLE_CONNS=25
ENV CREST_DB_MAX_IDLE_TIME=15m

EXPOSE ${CREST_PORT}

ENTRYPOINT crest-app -port=${CREST_PORT} -env=${CREST_ENV} -dsn-path=${CREST_DSN_PATH} \
-db-max-open-conns=${CREST_DB_MAX_OPEN_CONNS} -db-max-idle-conns=${CREST_DB_MAX_OPEN_CONNS} \
-db-max-idle-time=${CREST_DB_MAX_IDLE_TIME}

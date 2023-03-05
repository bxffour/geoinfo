# syntax=docker/dockerfile:1

## 
## BUILD
##
FROM golang:1.20.1-alpine3.17 as build-env

COPY . /go/src/app
WORKDIR /go/src/app

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /go/bin/crest-app ./cmd/api

##
## DEPLOY
##
FROM alpine:3.17.2
COPY --from=build-env /go/bin/crest-app /usr/local/bin/

ENV CREST_PORT=8080

EXPOSE ${CREST_PORT}

CMD crest-app --port=${CREST_PORT} --env=${CREST_ENV} --dsn-path=${CREST_DSN_PATH} \
--db-max-open-conns=${CREST_DB_MAX_OPEN_CONNS} --db-max-idle-conns=${CREST_DB_MAX_OPEN_CONNS} \
--db-dsn=${CREST_DB_DSN} --db-max-idle-time=${CREST_DB_MAX_IDLE_TIME}

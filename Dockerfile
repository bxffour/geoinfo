# syntax=docker/dockerfile:1

## 
## BUILD
##
FROM golang:1.20.1-alpine3.17 as build-env

COPY . /go/src/app
WORKDIR /go/src/app

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/crest-app ./cmd/api

##
## DEPLOY
##
FROM alpine:3.17.2

ENV SERVICE_USER=geoinfo \
    SERVICE_UID=1001 \
    SERVICE_GROUP=geoinfo \
    SERVICE_GID=1001

RUN addgroup -g ${SERVICE_GID} ${SERVICE_GROUP} && \
    adduser -D -H -G ${SERVICE_GROUP} -s /sbin/nologin -u ${SERVICE_UID} ${SERVICE_USER}

COPY --from=build-env /go/bin/crest-app /usr/local/bin/

ENV CREST_PORT=8080

USER ${SERVICE_USER}

ARG IMAGE_VERSION
ARG IMAGE_REVISION
ARG IMAGE_CREATED

LABEL org.opencontainers.image.title="geoinfo" \
      org.opencontainers.image.description="A REST API for getting information about countries" \
      org.opencontainers.image.url="https://ghcr.io/bxffour/crest/api" \
      org.opencontainers.image.source="https://github.com/bxffour/crest-countries" \
      org.opencontainers.image.vendor="thi-startup" \
      org.opencontainers.image.version="$(IMAGE_VERSION)" \
      org.opencontainers.image.licenses="GPLv3" \
      org.opencontainers.image.authors="Nana Kwadwo <agyemangclinton8@gmail.com>" \
      org.opencontainers.image.created="$(IMAGE_CREATED)" \
      org.opencontainers.image.revision="$(IMAGE_REVISION)"

EXPOSE ${CREST_PORT}

CMD crest-app --port=${CREST_PORT} --db-max-open-conns=${CREST_DB_MAX_OPEN_CONNS} \
    --db-max-idle-conns=${CREST_DB_MAX_OPEN_CONNS} \
    --db-user=${CREST_DB_USER} --db-password=${CREST_DB_PASSWORD} --db-dbname=${CREST_DATABASE} \
    --db-port=${CREST_DB_PORT} --db-host=${CREST_DB_HOST} \
    --db-dsn=${CREST_DB_DSN} --db-max-idle-time=${CREST_DB_MAX_IDLE_TIME}

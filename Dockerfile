# syntax=docker/dockerfile:1

## 
## BUILD
##
FROM golang:1.20.1-alpine3.17 as build-env

COPY . /go/src/app
WORKDIR /go/src/app

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/geoinfo-api ./cmd/api

##
## DEPLOY
##
FROM alpine:3.17.2

RUN apk add --no-cache postgresql-client
RUN apk add --no-cache bash

ENV SERVICE_USER=geoinfo \
    SERVICE_UID=1001 \
    SERVICE_GROUP=geoinfo \
    SERVICE_GID=1001

RUN addgroup -g ${SERVICE_GID} ${SERVICE_GROUP} && \
    adduser -D -H -G ${SERVICE_GROUP} -s /sbin/nologin -u ${SERVICE_UID} ${SERVICE_USER}

COPY --from=build-env /go/bin/geoinfo-api /usr/local/bin/
COPY ./geoinfo-start.sh /bin/gstart 
RUN chmod +x /bin/gstart

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

CMD [ "gstart", "geoinfo-api"]

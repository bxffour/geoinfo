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
EXPOSE 8080
CMD [ "crest-app" ]

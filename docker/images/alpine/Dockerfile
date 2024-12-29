ARG GO_DOCKER_VERSION=1.23

FROM golang:${GO_DOCKER_VERSION}-alpine AS builder

ENV APP_SRC=/go/src \
    APP_Bin=/go/bin \
    APP_ENV=development \
    APP_LISTEN_PORT=8080

COPY . ${APP_SRC}

RUN go build -o bin ${APP_SRC}/cmd/app/main.go

CMD ["tail", "-f", "/dev/null"]
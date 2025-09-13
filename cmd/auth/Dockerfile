FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /var/backend

COPY cmd/auth/main.go main.go
COPY . .

RUN go mod tidy
RUN go build -o main main.go

FROM alpine:edge AS prod

RUN apk add bash

COPY --from=build /var/backend/main /app/main
COPY --from=build /var/backend/.env /app/.env

RUN mkdir -p /app/files
COPY files /app/files

RUN mkdir -p /app/configs
COPY configs app/configs

WORKDIR /app
EXPOSE 44000

ENTRYPOINT ["./main"]
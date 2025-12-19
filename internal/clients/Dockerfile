FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /var/backend

COPY internal/clients/. clients/
COPY . .

RUN go mod download
RUN go mod tidy
RUN go build -o main clients/main.go clients/auth.go clients/relations.go clients/profiles.go clients/chats.go clients/kafka.go

FROM alpine:edge AS prod

RUN apk add bash

COPY --from=build /var/backend/main /app/main
COPY --from=build /var/backend/.env /app/.env

RUN mkdir -p /app/files
COPY files /app/files   

RUN mkdir -p /app/configs
COPY configs app/configs

WORKDIR /app
EXPOSE 8000

ENTRYPOINT ["./main"]
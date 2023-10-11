FROM golang:1.21.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk update && apk add netcat-openbsd
WORKDIR /app
COPY --from=builder /app/main .
COPY wait-for.sh .
COPY template /app/template

CMD [ "/app/main" ]
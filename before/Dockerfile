# Builder
FROM golang:1.13.0-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o before main.go

# Distribution
FROM alpine:latest
RUN apk add tzdata
COPY --from=builder app/before /app/before
COPY --from=builder app/config.yaml /app/config.yaml
WORKDIR /app
EXPOSE 9595
CMD ./before
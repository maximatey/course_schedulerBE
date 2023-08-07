# Build stage
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git gcc libc-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src/ ./src/
# RUN go build -o main ./src/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./src/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -o /main
# Final stage
FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
LABEL Name=courseschedulerbe Version=0.0.1

CMD [ "./main" ]

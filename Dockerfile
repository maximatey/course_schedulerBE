# Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# COPY *.go ./
# RUN go build -o main src/main.go
COPY src/ ./src/
RUN go build -o main ./src/main.go

# RUN CGO_ENABLED=0 GOOS=linux go build -o /main
# Final stage
FROM alpine:latest
COPY --from=builder /app .
EXPOSE 8080
LABEL Name=courseschedulerbe Version=0.0.1

CMD [ "./main" ]

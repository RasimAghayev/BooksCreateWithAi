FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o uploader.exe ./cmd/uploader

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/uploader.exe .
COPY --from=builder /app/.env.example ./
EXPOSE 8080
CMD ["./uploader.exe"]

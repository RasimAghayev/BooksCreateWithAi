FROM golang:1.24.5-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd

FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/app
COPY --from=build /app/internal/migrations ./internal/migrations
CMD ["./app"]

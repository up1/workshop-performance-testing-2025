FROM golang:1.24.4-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go mod tidy
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -o api main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
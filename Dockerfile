FROM mirror.gcr.io/library/golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/server ./cmd/server

FROM mirror.gcr.io/library/alpine
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /bin/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8000
ENTRYPOINT ["./server"]
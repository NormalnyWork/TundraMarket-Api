FROM mirror.gcr.io/library/golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go build -o /bin/server ./cmd/server

FROM mirror.gcr.io/library/alpine:3.23
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /bin/server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8000
ENTRYPOINT ["./entrypoint.sh"]
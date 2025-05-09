FROM golang:1.23.4 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/app/main ./cmd/app

FROM alpine AS app

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/cmd/app/main /app/scheduler.db ./
COPY --from=builder  /app/web ./web

ENV TODO_DBFILE=scheduler.db
ENV TODO_WEB_DIR=./web/

CMD ["./main"]
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./server -ldflags "-w -s" ./cmd/server/

FROM alpine

WORKDIR /app

COPY --from=builder /app/server ./server

EXPOSE 8080

ENTRYPOINT [ "./server" ]

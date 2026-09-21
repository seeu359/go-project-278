FROM golang:1.26-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest && \
    go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

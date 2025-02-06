FROM golang:1.23.5-alpine

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]

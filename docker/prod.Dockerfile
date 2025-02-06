FROM golang:1.23.5-alpine as builder

RUN apk update && apk upgrade && \
  apk --update add git make bash build-base

WORKDIR /app

COPY . .

# Build app executable
RUN make build

# Build db migrations executable
RUN make build-migrate

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
  apk --update --no-cache add tzdata && \
  mkdir /app 

WORKDIR /app 

EXPOSE 9090

COPY --from=builder /app/bin/engine /app/
COPY --from=builder /app/bin/migrate /app/

CMD /app/engine

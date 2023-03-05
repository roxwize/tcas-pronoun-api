# use official go runtime
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . /app

# download dependencies
RUN go mod download

# build the web server binary
RUN go build -o server

EXPOSE 1337

CMD ["./server"]

FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o app ./examples/main.go

CMD ["./app"]

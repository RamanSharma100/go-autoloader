FROM golang:1.26-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go test ./... -v

ENV CGO_ENABLED=1 GOOS=linux

RUN go build -o /app/bin/app ./examples/main.go

CMD ["/app/bin/app"]
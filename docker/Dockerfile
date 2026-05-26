FROM golang:1.26.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY 'cmd' './cmd'
COPY internal ./internal

RUN go build -o /usr/local/bin/lifeline ./cmd/app/main.go

CMD ["lifeline"]

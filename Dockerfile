FROM golang:1.21.0

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o effective_mobile ./cmd/effective_mobile/main.go

EXPOSE 8080

CMD ["./effective_mobile"]
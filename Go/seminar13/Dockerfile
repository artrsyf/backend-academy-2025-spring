FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN go mod download && go build -o app cmd/app/main.go

EXPOSE 8080

CMD ["./app"]
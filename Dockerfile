FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN HTTPS_PROXY=http://fodev.org:8118 go build -o app main.go

CMD ["./app"]
FROM golang:1.18-alpine

WORKDIR /app
COPY . .

RUN go build -o go-gin

EXPOSE 8080

CMD ./go-gin
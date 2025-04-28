FROM golang:1.24-alpine

RUN apk add --no-cache git

RUN go version

WORKDIR /app

COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/app/main.go"]
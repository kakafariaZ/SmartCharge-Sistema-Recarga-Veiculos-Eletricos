FROM golang:1.24.1

WORKDIR /app

COPY . .

RUN go build -o server

EXPOSE 8080

CMD ["./server"]


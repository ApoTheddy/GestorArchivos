FROM golang:1.16

WORKDIR /app

COPY . .

RUN go build -o gfilesapp

EXPOSE 3000

CMD ["./gfilesapp"]

FROM golang

WORKDIR /app

COPY . .

RUN go build -o jobs-runner cmd/jobs/main.go

CMD ["./jobs-runner"]
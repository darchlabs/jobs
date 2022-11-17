FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/jobs cmd/jobs/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin/jobs bin/jobs
ENTRYPOINT [ "./bin/jobs" ]
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./task .


FROM alpine:3

COPY --from=builder /app/task ./
RUN chmod +x ./task

ENTRYPOINT ["./task"]
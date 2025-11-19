FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./task .


FROM alpine:3

COPY --from=builder /app/task ./
RUN chmod +x ./task

ENV WORKLOAD_SIZE=20000
ENV TASK_ID="task-id"
ENV EXECUTION_SITE="local"
ENV DEVICE_ID=0
ENV CALLBACK_ADDR="http://localhost:8080"

ENTRYPOINT ["./task"]
FROM golang:1.14.2 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR ./src/github.com/firefly-crm/fireflycrm-bot-backend
ADD . .
RUN go mod download
RUN go build -o main cmd/app/app.go

FROM alpine
WORKDIR /app

COPY --from=builder /go/src/github.com/firefly-crm/fireflycrm-bot-backend/main .
CMD ["./main"]

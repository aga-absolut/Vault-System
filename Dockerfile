FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /server ./cmd/server/main.go
RUN go build -o /client ./cmd/client/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /server .
COPY --from=builder /client .
COPY secret.json .
EXPOSE 50051          
CMD ["./server"]
FROM golang:1.17-alpine as builder

WORKDIR /app

COPY cmd/ cmd/

RUN go build -o hello-world-app cmd/hello-world-app/*

# Final image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/hello-world-app .

EXPOSE 8080

CMD ["./hello-world-app"]

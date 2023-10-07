# build stage
FROM golang:1.20-alpine AS builder
WORKDIR /go/src/github.com/twibber/api
COPY . .
RUN go mod download && go build -ldflags="-s -w" -o build/app main.go

# deploy stage
FROM alpine:latest
COPY --from=builder /go/src/github.com/twibber/api/build/app /app
COPY --from=builder /go/src/github.com/twibber/api/.env* /
COPY --from=builder /go/src/github.com/twibber/api/mailer/templates/ /mailer/templates/
EXPOSE 8080
ENTRYPOINT /app
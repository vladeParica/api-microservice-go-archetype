# Build stage
FROM golang:1.20-alpine AS builder

LABEL maintainer="vparica"

ENV GO111MODULE=on
WORKDIR /go-microservice

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /api-microservice-go-archetype ./cmd/main.go


# Final stage
FROM alpine:latest

WORKDIR "/microservice-app"

COPY --from=builder /api-microservice-go-archetype .

EXPOSE 8000

RUN chmod +x ./go-microservice-template

ENTRYPOINT ["/microservice-app/api-microservice-go-archetype"]

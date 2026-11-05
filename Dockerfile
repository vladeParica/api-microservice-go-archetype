# syntax=docker/dockerfile:1

ARG GO_VERSION=1.27

# ---- Build stage ----
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/api ./cmd/api

# ---- Runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/api /app/api
COPY configs/config.yaml /app/configs/config.yaml

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/api"]

FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm@sha256:ee28302232bca53c6cfacf0b00a427ebbda10b33731c78d3dcf9f59251b23c9c AS base

WORKDIR /workspaces/furina

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /workspaces/furina/bin/app ./cmd

FROM base AS development
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

FROM alpine:latest@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c AS production
COPY --from=builder /workspaces/furina/bin/app .
CMD ["/app"]

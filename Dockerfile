FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm AS base

WORKDIR /workspaces/mimo

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

FROM base AS development
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS production
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /workspaces/mimo/bin/app ./cmd
CMD ["/workspaces/mimo/bin/app"]

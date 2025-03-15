FROM golang:1.24 AS base

WORKDIR /usr/src/mimo

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

FROM base AS development
CMD ["go", "run", "./..."]

FROM base AS production
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/mimo ./cmd
CMD ["/usr/local/bin/mimo"]

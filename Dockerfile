FROM golang:1.24 AS base

WORKDIR /usr/src/mimo

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM alpine:latest
WORKDIR /usr/src/mimo
COPY --from=builder /usr/src/mimo /usr/src/mimo

EXPOSE 8080

FROM base AS dev
CMD ["go", "run", "./..."]

FROM base AS prod
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/mimo ./cmd
CMD ["/usr/local/bin/mimo"]

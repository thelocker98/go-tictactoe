# Build stage
FROM golang:tip-alpine3.22 AS builder

# Install build deps for go-sqlite3
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY . .

# Enable CGO and build
ENV CGO_ENABLED=1
RUN go build -o /tictactoe main.go

# Final stage
FROM alpine:3.22

# Install runtime SQLite libraries only
RUN apk add --no-cache sqlite-libs

WORKDIR /srv
COPY --from=builder /tictactoe /srv/tictactoe
COPY --from=builder /app/templates /srv/templates

EXPOSE 8080
CMD ["/srv/tictactoe"]

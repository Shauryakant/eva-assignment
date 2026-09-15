# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install ca-certificates in builder stage if needed
RUN apk add --no-cache ca-certificates git

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build lightweight statically compiled binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o ticket-system main.go

# Production Stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/ticket-system .

EXPOSE 8080

ENTRYPOINT ["/app/ticket-system"]

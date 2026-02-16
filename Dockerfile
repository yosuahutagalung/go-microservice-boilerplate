# --- STAGE 1: Builder ---
# We use the matching Debian-based Go image for consistency
FROM golang:1.22-bookworm AS builder

# Install make and git
RUN apt-get update && apt-get install -y make git

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/bin/server ./cmd/service_boilerplate/

# --- STAGE 2: Runner ---
FROM debian:bookworm-slim AS runner

# Install CA certificates (for HTTPS), tzdata (for timezones), and curl (for debugging)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set the default timezone
ENV TZ=Asia/Jakarta

WORKDIR /app

# Copy the compiled binary and configs from the builder
COPY --from=builder /app/bin/server ./server
COPY --from=builder /src/configs ./configs

# Expose HTTP and gRPC ports
EXPOSE 8000
EXPOSE 9000

# Start the application
CMD ["./server", "-conf", "./configs"]

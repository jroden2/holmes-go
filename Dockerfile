# syntax=docker/dockerfile:1

# --- STAGE 1: Builder ---
FROM golang:1.25 AS builder
WORKDIR /src/app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the static binary
# We output it to the current directory to keep paths simple
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o holmes ./cmd/main.go


# --- STAGE 2: Final Runtime ---
FROM gcr.io/distroless/static:nonroot

# Set the working directory - your app will look for ./templates here
WORKDIR /app

# Documentation for the exposed port
EXPOSE 8080

# Copy the binary and the templates folder from the builder
COPY --from=builder /src/app/holmes .
COPY --from=builder /src/app/templates ./templates
USER nonroot:nonroot

ENTRYPOINT ["./holmes"]
# syntax=docker/dockerfile:1
FROM golang:1.25 AS builder
WORKDIR /src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o holmes ./cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app

EXPOSE 8080

COPY --from=builder /src/app/holmes .
COPY --from=builder /src/app/templates ./templates
USER nonroot:nonroot

ENTRYPOINT ["./holmes"]
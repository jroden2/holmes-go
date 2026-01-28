# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.25 AS builder
WORKDIR /src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o holmes ./cmd/main.go

FROM --platform=$BUILDPLATFORM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /src/app/holmes .
COPY --from=builder /src/app/templates ./templates
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["./holmes"]

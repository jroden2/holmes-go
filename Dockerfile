# syntax=docker/dockerfile:1
FROM gcr.io/distroless/static:nonroot
WORKDIR /app

EXPOSE 8080

COPY holmes .
COPY templates ./templates

USER nonroot:nonroot
ENTRYPOINT ["./holmes"]

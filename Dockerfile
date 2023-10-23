FROM golang:1.21 AS builder

WORKDIR /app
COPY *.go go.mod go.sum /app/
RUN GGO_ENABLED=0 GOOS=linux go build -o main

FROM debian:buster-slim AS runner
COPY --from=builder /app/main /main
USER nobody
CMD ["/main"]

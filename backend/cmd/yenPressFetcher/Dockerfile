FROM registry.suse.com/bci/golang:1.25 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .

RUN CGO_ENABLED=0 GOOS=linux go build -o /yenPressFetcher ./cmd/yenPressFetcher

FROM registry.suse.com/bci/bci-micro:latest

WORKDIR /app

COPY --from=builder /yenPressFetcher /usr/local/bin/yenPressFetcher

ENTRYPOINT ["yenPressFetcher"]

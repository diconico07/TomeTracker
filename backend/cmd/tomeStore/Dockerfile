# Stage 1: Build the Go application
FROM registry.suse.com/bci/golang:1.25 AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .

RUN CGO_ENABLED=0 GOOS=linux go build -o /tomeStore ./cmd/tomeStore

# Stage 1.5: Build the Vue frontend
FROM registry.suse.com/bci/nodejs:22 AS frontend_builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

COPY frontend/ .
RUN npm run build


# Stage 2: Create the final, minimal image
FROM registry.suse.com/bci/bci-micro:latest

WORKDIR /app

COPY --from=builder /tomeStore /usr/local/bin/tomeStore
COPY --from=frontend_builder /app/dist ./dist 

EXPOSE 8080

ENTRYPOINT ["tomeStore"]

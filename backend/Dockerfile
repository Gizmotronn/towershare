FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pocketbase ./main.go

# ---

FROM alpine:3.19

RUN apk add --no-cache ca-certificates wget

WORKDIR /pb

COPY --from=builder /build/pocketbase .

EXPOSE 8090

ENTRYPOINT ["/pb/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/pb/pb_data", "--migrationsDir=/pb/pb_migrations", "--hooksDir=/pb/pb_hooks"]

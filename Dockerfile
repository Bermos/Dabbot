# Build dabbot
FROM golang:1.17.1-alpine AS backend

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY Main.go Main.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dabbot .

# Deploy stage
FROM alpine

WORKDIR /app
COPY --from=backend /build/dabbot ./dabbot

ENTRYPOINT ["./dabbot"]
# Build dabbot
FROM golang:1.18-alpine AS backend

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go main.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dabbot .

# Deploy stage
FROM alpine

WORKDIR /app
COPY --from=backend /build/dabbot ./dabbot

VOLUME /app/dabs

ENTRYPOINT ["./dabbot"]

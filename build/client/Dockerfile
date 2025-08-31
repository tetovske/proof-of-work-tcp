# Use the official Golang image to create a build artifact.
FROM golang:1.24 as builder

# Create a working directory.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod ./

# Download dependencies.
RUN go mod download

# Copy the source code into th container.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -o client ./cmd/client

# Use a small image to run the server
FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/client .

CMD ["./client"]

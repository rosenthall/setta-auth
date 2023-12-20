# Dockerfile for gRPC service (JWT microservice)
FROM golang:1.16 as builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the application using the Makefile
RUN make -C build/ run

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/build/auth .

# Start the application
CMD ["./target/auth"]

# Dockerfile for gRPC service (JWT microservice)
FROM golang:1.21.1 as builder

# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod download

# Build the application via Makefile
RUN make -C build/ build

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Create a directory for the app
RUN mkdir /root/app

# Copy the built application from the builder stage
COPY --from=builder /app/build/target/auth /root/app/auth

# Copy configuration file and keys from the project
COPY config.toml /root/config.toml
COPY keys/ /root/keys/

# Set the working directory to /root/app
WORKDIR /root/app

# Start the application
CMD ["./auth"]

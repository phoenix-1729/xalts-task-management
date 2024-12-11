# Use the official Golang image as a builder
FROM golang:1.19 AS builder

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Give appropriate permissions to the app directory and its contents
RUN chmod -R 777 /app

# Build the application
RUN go build -o main .

# Use a lightweight image for production
FROM ubuntu:latest

# Install necessary dependencies for running Go binaries
RUN apt-get update && apt-get install -y ca-certificates

# Set the working directory for the second stage
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
RUN ls -l /app
# Make the binary executable
RUN chmod +x /app/main

# Expose port 8002
EXPOSE 8002

# Command to run the application
CMD ["./main"]

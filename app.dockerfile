# Use a Golang image for building the application
FROM golang:1.23-alpine AS build

# Install necessary dependencies
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY gRPC-Account-service/go.mod gRPC-Account-service/go.sum ./

# Download the Go modules dependencies (optional, but useful for syncing)
RUN GO111MODULE=on go mod download

# Ensure the dependencies are in sync
RUN GO111MODULE=on go mod tidy

# Copy the source code
COPY gRPC-Account-service/ ./

# Build the application binary
RUN go build -o /app/bin/account ./cmd/account

# Final stage: create a smaller image to run the application
FROM alpine:3.11

# Set the working directory inside the container
WORKDIR /usr/bin

# Copy the compiled binary from the build stage
COPY --from=build /app/bin/account .

# Expose the port your application runs on
EXPOSE 8080

# Command to run the application
CMD ["./account"]


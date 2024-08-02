# Build stage
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Define build arguments for build tags
ARG BUILD_TAGS="ttl persistent"

# Build the Go application with the desired build tags
RUN CGO_ENABLED=0 GOOS=linux go build -tags "$BUILD_TAGS" -o justkv .

# Final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/justkv .

# Set environment variable for the server port
ENV PORT 8080

# Expose the port on which the server will run
EXPOSE 8080

# Command to run the binary
CMD ["./justkv"]

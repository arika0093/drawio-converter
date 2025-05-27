# Use an official lightweight image as a parent image for building
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o drawio-converter .

# Use a minimal image for running the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/drawio-converter .

# Expose the default port for the server
EXPOSE 8080

# Set the ENTRYPOINT to the application executable
ENTRYPOINT ["/app/drawio-converter"]

# Use CMD to specify default arguments
CMD ["--server"]

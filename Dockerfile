# Stage 1: Build the Frontend
FROM node:18.14.1 AS frontend

# Set the working directory for the frontend
WORKDIR /app

# Copy the entire project
COPY . /app

# Install make, required for the build process
RUN apt-get update && apt-get install -y make

# Change to the ./web directory and build the frontend
WORKDIR /app
RUN make build-web

# Stage 2: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory for Go application
WORKDIR /src

# Copy Go modules files
COPY go.sum go.mod ./

# Download Go modules
RUN go mod download

# Copy the entire project
COPY . .

# Move frontend files to the root
COPY --from=frontend /app/web ./web

# Build the Go application
RUN CGO_ENABLED=0 go build -o /bin/app .

# Stage 3: Create the final image
FROM ubuntu:latest

# Install dependencies for the final image
RUN apt-get update && apt-get -y upgrade && apt-get install -y --no-install-recommends \
  libssl-dev \
  ca-certificates \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Copy the Go binary from the builder stage
COPY --from=builder /bin/app /checkpointz

# Expose port 5555
EXPOSE 5555

# Set the entrypoint for the container
ENTRYPOINT ["/checkpointz"]

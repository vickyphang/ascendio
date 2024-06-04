# Builder
FROM golang:1.21-alpine3.19 AS go

MAINTAINER Vicky Phang <vickyphang11@gmail.com>

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY main.go ./

# Build
RUN go build -o ascendio


# Actual image
FROM alpine:3.14

# Set working directory
WORKDIR /app

# Copy from builder
COPY --from=go /app/ascendio /app/ascendio

# Ports the application is going to listen on by default
EXPOSE 8080

# Run
CMD ["/app/ascendio"]
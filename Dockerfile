FROM golang:1.21-alpine3.19

MAINTAINER Vicky Phang <vickyphang11@gmail.com>

# Set destination for COPY
WORKDIR /app


# Copy the source code
COPY go.mod go.sum main.go ./

# Download Go modules
RUN go mod download

# Build
RUN go build -o ./ascendio

# Ports the application is going to listen on by default
EXPOSE 8080

# Run
CMD ["./ascendio"]
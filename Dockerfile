# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.12 as builder

# Add Maintainer Info
LABEL maintainer="dilap54 <dilap54@mail.ru>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o idor ./cmd


######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/idor ./
COPY  web ./web

# Command to run the executable
CMD ["sh","-c","/root/idor"]
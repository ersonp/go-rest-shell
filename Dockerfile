# Start from the latest golang base image
FROM golang:latest

COPY . /go-rest-shell

# Set the Current Working Directory inside the container
WORKDIR /go-rest-shell

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o rest-shell ./cmd

# Command to run the executable
CMD ["./rest-shell", "-host", "0.0.0.0"]
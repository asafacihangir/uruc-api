# Start from a Debian-based image with the Go 1.x SDK installed
FROM golang:1.19

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -v -o server


# Run the web service on container startup.
CMD ["/app/server"]
EXPOSE 8080

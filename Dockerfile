# Use the official Golang image as the base image
FROM golang:1.23.5-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o video_manager .

FROM scratch AS runtime

VOLUME videos

# Copy the binary from the build stage to the runtime stage
COPY --from=builder /app/video_manager /

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD [ "/video_manager", "web", "-v", "/videos"]
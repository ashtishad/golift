# Use a multi-stage build to keep the final image as small as possible.
# The builder stage uses the official Go image to compile the application.
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files(if exists), then download dependencies.
# Taking advantage of Docker's cache layers, only re-download dependencies if these files change.
# COPY go.mod go.sum ./
# RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the Go app.
# -o specifies the output binary name, here it's set to main.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Start a new stage to keep the final image small.
FROM alpine:latest

# Add ca-certificates in case you need to make HTTPS requests.
RUN apk --no-cache add ca-certificates

# Create a non-root user and group to run the application.
# This is a security best practice to limit the privileges of the application process.
RUN addgroup -S appgroup && adduser -S ash -G appgroup

# Set the working directory in the container.
WORKDIR /app

# Copy the binary from the builder stage to the current stage.
COPY --from=builder /app/main .

# Change to the non-root user.
USER ash

# Run the binary.
ENTRYPOINT ["./main"]

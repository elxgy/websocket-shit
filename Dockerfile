# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Compile the application.
# CGO_ENABLED=0 is important for creating a static binary.
# GOOS=linux ensures we build for a Linux environment.
RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-server .

# ---

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /root/

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /chat-server ./

# Change ownership to non-root user
RUN chown appuser:appgroup /root/chat-server

# Switch to non-root user
USER appuser

# Tell Docker that the container listens on port 8080
EXPOSE 8080

# The command to run when the container starts
ENTRYPOINT ["./chat-server"]

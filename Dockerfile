# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

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
FROM scratch

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /chat-server /chat-server

# Tell Docker that the container listens on port 8080
EXPOSE 8080

# The command to run when the container starts
ENTRYPOINT ["/chat-server"]

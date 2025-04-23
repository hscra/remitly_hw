# Define base image
FROM golang:1.24

# Working directory inside the container
WORKDIR /app

# Copy Go mod and sum files
COPY go.mod go.sum ./

# Install go modules into working directory inside the image. 
RUN go mod download 

# Download dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Expose port 8080 for the Go API
EXPOSE 8080

# Run the Go application
CMD ["go", "run", "main.go"]
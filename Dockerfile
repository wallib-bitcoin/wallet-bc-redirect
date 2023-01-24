# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go script
RUN go build -o main .

# Expose the port that the script binds to
EXPOSE 8080

# Run the script when the container starts
CMD ["./main"]

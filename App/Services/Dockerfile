# Use a minimal base image.
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the previous stage (build stage)
COPY helloApp /app/helloApp

# Expose the port on which your application listens
EXPOSE 8080

# Define the command to run your application
CMD ["./helloApp"]


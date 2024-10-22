# syntax=docker/dockerfile:1

# Step 1: Build the Go application
FROM golang:1.22-alpine3.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Install dependencies
RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /Simulation_Modbus

# Step 2: Create the final lightweight image
FROM alpine:3.18

# Copy the built binary from the builder stage
COPY --from=builder /Simulation_Modbus /

# Copy the YAML configuration files into the container
COPY Sensor01.yml /Sensor01.yml
COPY Sensor02.yml /Sensor02.yml

# Expose the port your application listens on (adjust if needed)
EXPOSE 502

# Set the command to run your application
ENTRYPOINT  ["/Simulation_Modbus"]
CMD ["Sensor01.yml"]

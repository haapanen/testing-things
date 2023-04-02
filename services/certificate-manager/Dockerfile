# Use the official golang image as the builder
FROM golang:1.19 as builder

# Set the working directory within the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the source code to the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o certificate-manager .

# Start a new stage using the scratch image
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/certificate-manager /certificate-manager

# Set the entrypoint for the container
ENTRYPOINT ["/certificate-manager"]

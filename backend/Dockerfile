FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build
# Copy the code into the container
COPY ./ ./
#COPY go.mod go.sum ./

# Copy and download dependency using go mod
RUN go mod download

# Build the application
RUN go build -ldflags '-w -s' -a -o ./serverBuild ./

# Build a small image
FROM alpine

#create local directory for workspace
WORKDIR /app/workspace
# Move to /app directory
WORKDIR /app
#no need to copy .env, env variables will be loaded from docker-compose
#COPY .env .
# Copy binary from build to main folder
COPY --from=builder /build/serverBuild ./server
# Export necessary port
EXPOSE 8080

# Command to run
#CMD ["sleep", "infinity"]
ENTRYPOINT ["/app/server", "frontend"]

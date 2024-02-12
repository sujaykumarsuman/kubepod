# Build stage
FROM golang:1.21.0-alpine3.18 as builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

# Build the Go application and set file permissions
RUN CGO_ENABLED=0 GOOS=linux go build -o /kubepod cmd/main.go && chmod +x /kubepod

# Final stage
FROM alpine:3.18

# Copy the binary from the builder stage
COPY --from=builder /kubepod /kubepod

# Set the default command
CMD ["/kubepod"]

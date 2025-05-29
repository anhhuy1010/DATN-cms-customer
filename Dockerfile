FROM golang:1.24-alpine

WORKDIR /app

# Install necessary tools
RUN apk add --no-cache bash git protoc curl build-base

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy entire source
COPY . ./

# Generate swagger docs (optional)
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init ./controllers/* || true

# Build the Go app
RUN go build -o /main ./main.go

# Run
EXPOSE 8080
ENTRYPOINT ["/main"]

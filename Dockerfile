FROM golang:1.24-alpine

WORKDIR /app

# Cài các công cụ cần thiết
RUN apk add --no-cache bash git curl build-base protobuf

# Cài plugin protoc cho Go và gateway
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest

# Thêm GOPATH/bin vào PATH để các plugin protoc hoạt động
ENV PATH="/go/bin:${PATH}"

# Copy go.mod và tải module
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source
COPY . .

# Sinh swagger docs (tùy chọn)
RUN swag init ./controllers/* || true

# Build ứng dụng Go
RUN go build -o /main ./main.go

# Lắng nghe cổng HTTP
EXPOSE 8080

# Chạy app
ENTRYPOINT ["/main"]

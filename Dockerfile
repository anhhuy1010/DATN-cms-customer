FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache bash git curl build-base protobuf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest \
    && apk add --no-cache envoy

ENV PATH="/go/bin:${PATH}"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init ./controllers/* || true

RUN go build -o /main ./main.go

COPY envoy.yaml /app/envoy.yaml
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 8080 9003

ENTRYPOINT ["/app/entrypoint.sh"]

.PHONY: build-app start restart logs ssh-app swagger proto-user proto-users proto-customer build-web release-web

build-app:
	docker-compose build app
start:
	docker-compose up
restart:
	docker-compose restart
logs:
	docker logs -f DATN-cms-customer
ssh-app:
	docker exec -it DATN-cms-customer bash
swagger:
	swag init ./controllers/*
proto-user:
	protoc -I grpc/proto/user/ \
		-I /usr/include \
		--go_out=paths=source_relative,plugins=grpc:grpc/proto/user/ \
		grpc/proto/user/user.proto
proto-users:
	protoc -I grpc/proto/users/ \
		-I /usr/include \
		--go_out=paths=source_relative,plugins=grpc:grpc/proto/users/ \
		grpc/proto/users/users.proto
proto-customer:
	protoc -I grpc/proto/customer \
       -I third_party/googleapis \
       -I third_party/protobuf/src \
       --go_out=paths=source_relative:grpc/proto/customer \
       --go-grpc_out=paths=source_relative:grpc/proto/customer \
       --grpc-gateway_out=paths=source_relative:grpc/proto/customer \
       grpc/proto/customer/customer.proto

build-web:
	heroku container:push web -a imatching
release-web:
	heroku container:release web -a imatching

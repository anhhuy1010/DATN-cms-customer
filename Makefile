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

build-app:
	docker-compose build app
start:
	docker-compose up
restart:
	docker-compose restart
logs:
	docker logs -f DATN-cms-ideas
ssh-app:
	docker exec -it DATN-cms-ideas bash
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
proto-idea:
	protoc -I grpc/proto/idea/ \
		-I /usr/include \
		--go_out=paths=source_relative,plugins=grpc:grpc/proto/idea/ \
		grpc/proto/idea/idea.proto
push-app:
	heroku container:push web -a cms-ideas-app
build-app:
	heroku container:release web -a cms-ideas-app
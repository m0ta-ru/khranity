.PHONY: build-svc build-con

build-svc:
	go build -o ./build/service cmd/service/main.go

build-con:
	go build -o ./build/console cmd/console/*.go

.PHONY: image-svc image-svc-alpine docker-svc-con
# build service image
image-svc-alpine:
	docker build 									\
		--tag khranity:service .					\
		--no-cache									\
		--force-rm									\
		--build-arg GO_VERSION=1.20.6				\
		--file ./cmd/service/Dockerfile-alpine

image-svc-ubuntu:
	docker build 									\
		--tag khranity:service .					\
		--no-cache									\
		--force-rm									\
		--build-arg GO_VERSION=1.20.6				\
		--file ./cmd/service/Dockerfile-ubuntu
image-svc:
	make image-svc-alpine

.PHONY: image-con image-con-alpine docker-run-con
# build console image
image-con-alpine:
	docker build 									\
		--tag khranity:console .					\
		--no-cache									\
		--force-rm									\
		--build-arg GO_VERSION=1.20.6				\
		--file ./cmd/console/Dockerfile-alpine

image-con:
	make image-con-alpine

# run image
docker-run-svc:
	docker stop khranity-svc
	docker rm khranity-svc
	docker run 													\
		--name khranity-svc										\
		--detach												\
		--restart always 										\
		--env-file ~/.khranity/.env								\
		--volume ~/.khranity/logs:/exec/logs					\
		--read-only --volume ~/.khranity/config:/exec/config	\
		--read-only --volume ~/:/exec/data						\
		--tmpfs /tmp											\
		khranity:service

# run image
docker-run-con:
	docker stop khranity-con
	docker rm khranity-con
	docker run 													\
		--name khranity-con										\
		--volume /root/temp:/exec/data							\
		--read-only --volume ~/.khranity/config:/exec/config	\
		--tmpfs /tmp											\
		khranity:console get -n "khranity" -p "/exec/data" -l "/exec/config/lore.yml"

image-delete:
	docker rm $(docker ps -q -f status=exited)
	docker image prune

test-get:
	go run ./cmd/console/*.go get -n "m0ta-blog-nextjs" -p "/root/temp"

test-put:
	go run ./cmd/console/*.go put -n "m0ta-blog-nextjs"
.PHONY: build-svc build-con

build-svc:
	go build -o ./build/service cmd/service/*.go

build-con:
	go build -o ./build/console cmd/console/*.go

.PHONY: image-svc
# build service image
image-svc:
	docker build 									\
		--tag khranity:service .					\
		--no-cache									\
		--force-rm									\
		--build-arg GO_VERSION=1.22.1				\
		--file ./docker/service/Dockerfile
	docker image prune --force --filter label=stage=temp

.PHONY: image-con
# build console image
image-con:
	docker build 									\
		--tag khranity:console .					\
		--no-cache									\
		--force-rm									\
		--build-arg GO_VERSION=1.22.1				\
		--file ./docker/console/Dockerfile
	docker image prune --force --filter label=stage=temp

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

test-get:
	go run ./cmd/console/*.go get -n "test" -p "/tmp/khranity"

test-put:
	go run ./cmd/console/*.go put -n "test"
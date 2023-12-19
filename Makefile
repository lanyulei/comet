PASSWORD=$(shell date +%s | base64 | head -c 32)
WORKDIR=$(shell pwd)
default:
	@echo "工作目录：$(shell pwd)"
postgres:
	mkdir -p data/postgres
	docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=${PASSWORD} -e PGDATA=/var/lib/postgresql/data/pgdata -v ${WORKDIR}/data/postgres:/var/lib/postgresql/data/pgdata postgres
	@echo "postgres 持久化数据目录: ${WORKDIR}/data/postgres，默认密码：${PASSWORD}"
redis:
	mkdir -p data/redis
	docker run -itd -p 6379:6379 --name redis -v ${WORKDIR}/data/redis:/data redis redis-server --requirepass ${PASSWORD}
	@echo "redis 持久化数据目录: ${WORKDIR}/data/redis，默认密码：${PASSWORD}"
build:
	@echo "go build ."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o comet main.go
migrate: build
	@echo "同步数据"

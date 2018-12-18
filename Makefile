include .env

# Run docker container with ClickHouse and create pmm DB.
up:
	mkdir logs
	docker-compose up -d --build
	sleep 60
	docker exec ch-server clickhouse client -h 127.0.0.1 --query="CREATE DATABASE IF NOT EXISTS pmm;"

# Remove docker containers.
down:
	docker-compose down --volumes
	rm -rf logs


pmm-up:
	docker-compose up -d


deploy:
	# docker exec pmm-server supervisorctl reload
	docker exec pmm-server supervisorctl stop qan-api2
	docker cp percona-qan-api2 pmm-server:/usr/sbin/percona-qan-api2
	docker exec pmm-server supervisorctl start qan-api2

# Connect to pmm DB.
ch-client:
	docker exec -ti ch-server clickhouse client -d pmm

ps-client:
	docker exec -ti ps-server mysql -uroot -psecret

# Run qan-api with envs.
# env $(cat .env | xargs) go run *.go
go-run:
	@echo "  > Runing with envs..." 
	go run *.go

# Pack ClickHouse migrations into go file.
go-generate:
	@echo "  >  Generating dependency files..."
	go-bindata -pkg migrations -o migrations/bindata.go -prefix migrations/sql migrations/sql
	protoc api/version/version.proto --go_out=plugins=grpc:.
	protoc api/collector/collector.proto --go_out=plugins=grpc:.

# Build binary.
linux-go-build: go-generate
	@echo "  >  Building binary..."
	GOOS=linux go build -o percona-qan-api2 *.go

# Build binary.
go-build: go-generate
	@echo "  >  Building binary..."
	go build -o percona-qan-api2 *.go

# Request API version.
# require https://github.com/uber/prototool
api-version:
	prototool grpc api/version --address 0.0.0.0:9911 --method version.Version/HandleVersion --data '{"name": "me"}'

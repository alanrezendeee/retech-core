.PHONY: run build docker push

run:
	go run ./cmd/api

build:
	go build -o bin/retech-core ./cmd/api

docker:
	docker build -f build/Dockerfile -t theretech/retech-core:0.1.0 .

push:
	docker push theretech/retech-core:0.1.0


build:
	docker build -t jacob-elektronik/rabbit-amazon-forwarder -f Dockerfile .

push: test build
	docker push jacob-elektronik/rabbit-amazon-forwarder

test:
	docker-compose run --rm tests

dev:
	go build

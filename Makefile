.PHONY: run build test clean docker-build docker-run

run:
	go run cmd/api/main.go

build:
	go build -ldflags="-s -w" -o qr-microservice cmd/api/main.go

test:
	go test ./...

clean:
	rm -f qr-microservice
	rm -f *.png

docker-build:
	docker build -t qr-go:latest .

docker-run:
	docker run --rm -p 8080:8080 qr-go:latest

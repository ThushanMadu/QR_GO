.PHONY: run build clean

run:
	go run cmd/api/main.go

build:
	go build -o qr-microservice cmd/api/main.go

clean:
	rm -f qr-microservice
	rm -f *.png

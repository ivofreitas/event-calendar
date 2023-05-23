install:
	go get ./...

start:
	go run main.go

test:
	go test -race ./...

cover:
	go test -cover ./...

configure:
	cp config/example.env config/.env
	swag init
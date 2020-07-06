build:
	go build ./cmd/auth-proxy.go

fmt:
	go fmt ./...

run: build
	./auth-proxy
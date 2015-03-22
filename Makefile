build:
	go build -o hazel-api api/cmd/main.go
	go build -o hazel-builder builder/cmd/main.go

run: build
	goreman start

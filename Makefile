BINARY_NAME=financialApp

build:
	go build -o ./bin/${BINARY_NAME} ./cmd/main.go

run: build
	./bin/${BINARY_NAME}

clean:
	go clean
	rm ./bin/${BINARY_NAME}

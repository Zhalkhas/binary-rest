BUILD_OUTPUT ?= ./build/app
build:
	go build -o $(BUILD_OUTPUT) cmd/app/main.go
run:
	go run cmd/app/main.go
test:
	go test  ./...
bench:
	go test -bench=. -count 10 -run=^# -benchtime=3s ./...

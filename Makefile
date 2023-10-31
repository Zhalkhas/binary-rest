build:
	go build -o ./build/app
run:
	go run cmd/app/main.go
test:
	go test  ./...
bench:
	go test -bench=. -count 10 -run=^# -benchtime=3s ./...

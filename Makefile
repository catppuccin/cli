cli: tidy
	go build -o build/ ./cmd/ctp
tidy:
	go mod tidy

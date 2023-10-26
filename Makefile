install: cli
	cp -v ctp ${HOME}/.local/bin
cli: tidy
	go build -o build/ ./cmd/ctp
tidy:
	go mod tidy

cli: tidy
  go build ./cmd/ctp -o build/
tidy:
  go mod tidy

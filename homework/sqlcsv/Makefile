build:
	go build -o sqlcsv main.go

fmt:
	gofmt -w -s .

test:
	go test -race -cover -coverprofile=coverage.out ./...

check:
	go vet ./...
	golangci-lint run ./...
test:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

bench:
	go test -v -benchmem -bench ./...

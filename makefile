test:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

bench:
	go test -v -benchmem -bench ./...

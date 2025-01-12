.PHONY: bench test


bench:
	go test -bench=. -benchmem

test:
	go test -v ./...

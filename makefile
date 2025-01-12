.PHONY: bench

bench: .bench

.bench:
	go test -bench=. -benchmem -count=5       


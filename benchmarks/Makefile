OUT := out
HEADER := $(OUT)/header
RESULTS := $(OUT)/results
xxh64 := $(OUT)/xxh64
wyhash := $(OUT)/wyhash
rapidhash := $(OUT)/rapidhash
rapidhash_comparable := $(OUT)/rapidhash.HashComparable
maphash_comparable := $(OUT)/maphash.Comparable

bench:
	@mkdir -p $(OUT) 2>/dev/null || true
	@rm -f $(OUT)/*

	@go test -count=10 -bench=. -benchmem | tee $(RESULTS)
	@grep -E '^(goos|goarch|pkg|cpu):' $(RESULTS) > $(HEADER)

	@cat $(HEADER) > $(xxh64)
	@grep '^BenchmarkXXH64/' $(RESULTS) | sed 's/BenchmarkXXH64\//Benchmark/' >> $(xxh64)
	@echo "PASS" >> $(xxh64)

	@cat $(HEADER) > $(wyhash)
	@grep '^BenchmarkWyhash/' $(RESULTS) | sed 's/BenchmarkWyhash\//Benchmark/' >> $(wyhash)
	@echo "PASS" >> $(wyhash)

	@cat $(HEADER) > $(rapidhash)
	@grep '^BenchmarkRapidhash/' $(RESULTS) | sed 's/BenchmarkRapidhash\//Benchmark/' >> $(rapidhash)
	@echo "PASS" >> $(rapidhash)

bench-comparable:
	@mkdir -p $(OUT) 2>/dev/null || true
	@rm -f $(OUT)/*

	@go test -tags=purego -count=10 -bench=BenchmarkComparable -benchmem | tee $(RESULTS)
	@grep -E '^(goos|goarch|pkg|cpu):' $(RESULTS) > $(HEADER)

	@cat $(HEADER) > $(rapidhash_comparable)
	@grep '^BenchmarkComparable/rapidhash/' $(RESULTS) | sed 's|rapidhash/||' >> $(rapidhash_comparable)
	@echo "PASS" >> $(rapidhash_comparable)

	@cat $(HEADER) > $(maphash_comparable)
	@grep '^BenchmarkComparable/maphash/' $(RESULTS) | sed 's|maphash/||' >> $(maphash_comparable)
	@echo "PASS" >> $(maphash_comparable)

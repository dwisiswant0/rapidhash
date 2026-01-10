PKG := rapidhash
OUT := out
HEADER := $(OUT)/header
RESULTS := $(OUT)/results
Hash := $(OUT)/Hash
Hasher := $(OUT)/Hasher
HashNano := $(OUT)/HashNano
HashMicro := $(OUT)/HashMicro

.PHONY: test
test:
	@go test -v -race .

.PHONY: build-test
build-test:
	@go test -c -o $(PKG).test .

.PHONY: bench
bench: build-test
	@mkdir -p $(OUT) 2>/dev/null || true
	@rm -f $(RESULTS)

	@./$(PKG).test -test.run - -test.count=10 -test.bench=. -test.benchmem -test.cpuprofile=cpu.out -test.memprofile=mem.out | tee $(RESULTS)
	@grep -E '^(goos|goarch|pkg|cpu):' $(RESULTS) > $(HEADER)

	@cat $(HEADER) > $(Hash)
	@grep '/Hash-' $(RESULTS) | sed 's|Computes/||' | sed 's|/Hash||' >> $(Hash)
	@echo "PASS" >> $(Hash)

	@cat $(HEADER) > $(Hasher)
	@grep '/Hasher-' $(RESULTS) | sed 's|Computes/||' | sed 's|/Hasher||' >> $(Hasher)
	@echo "PASS" >> $(Hasher)

	@cat $(HEADER) > $(HashNano)
	@grep '/HashNano-' $(RESULTS) | sed 's|Computes/||' | sed 's|/HashNano||' >> $(HashNano)
	@echo "PASS" >> $(HashNano)

	@cat $(HEADER) > $(HashMicro)
	@grep '/HashMicro-' $(RESULTS) | sed 's|Computes/||' | sed 's|/HashMicro||' >> $(HashMicro)
	@echo "PASS" >> $(HashMicro)

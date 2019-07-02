GO = $(shell which go)

bin/cretag: cretag.go go.mod go.sum
	$(GO) build -o $@ $<

.PHONY: clean
clean:
	$(RM) -r bin/

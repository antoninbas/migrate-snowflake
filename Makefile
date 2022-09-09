GO                 ?= go
BINDIR := $(CURDIR)/bin

all: bin

.PHONY: bin
bin:
	$(GO) build -o $(BINDIR)/migrate-snowflake github.com/antoninbas/migrate-snowflake

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: clean
clean:
	rm -rf bin

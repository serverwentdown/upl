GO = go
NPM = npm


.PHONY: all
all: build

.PHONY: clean
clean:
	$(RM) -r upl web/assets

.PHONY: build
build: upl

upl: *.go web/*.tmpl web
	$(GO) build -ldflags="-s -w" -tags "$(TAGS)" -v -o upl

.PHONY: test
test: web
	$(GO) test -cover -bench=. -v ./...
.PHONY: vet
vet: web
	$(GO) vet ./...


.PHONY: web
web: web/assets/bundle.js

web/node_modules:
	cd web && $(NPM) install

web/assets/bundle.js: web/node_modules
	cd web && $(NPM) run build

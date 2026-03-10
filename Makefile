BINARY := prayertime-cli
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3
GOVULNCHECK := go run golang.org/x/vuln/cmd/govulncheck@v1.1.4
GORELEASER := go run github.com/goreleaser/goreleaser/v2@v2.14.3

.PHONY: test lint vuln docs docs-check verify build release-check release-snapshot

test:
	go test ./...

lint:
	$(GOLANGCI_LINT) run ./...

vuln:
	$(GOVULNCHECK) ./...

docs:
	go run ./cmd/prayertime-cli-docs

docs-check:
	tmpdir="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmpdir"' EXIT; \
	go run ./cmd/prayertime-cli-docs \
		--docs-dir "$$tmpdir/docs/cli" \
		--man-dir "$$tmpdir/docs/man" \
		--completions-dir "$$tmpdir/completions" \
		--llms-full-path "$$tmpdir/llms-full.txt"; \
	diff -ru "$$tmpdir/docs/cli" docs/cli; \
	diff -ru "$$tmpdir/docs/man" docs/man; \
	diff -ru "$$tmpdir/completions" completions; \
	diff -u "$$tmpdir/llms-full.txt" llms-full.txt

verify: test lint vuln docs-check build

build:
	go build ./cmd/prayertime-cli

release-check:
	$(GORELEASER) check

release-snapshot:
	$(GORELEASER) release --snapshot --clean

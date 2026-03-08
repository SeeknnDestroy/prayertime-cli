BINARY := prayertime-cli

.PHONY: test lint vuln docs docs-check build release-snapshot

test:
	go test ./...

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

docs:
	go run ./cmd/prayertime-cli-docs

docs-check: docs
	git diff --exit-code -- docs/cli docs/man completions llms-full.txt

build:
	go build ./cmd/prayertime-cli

release-snapshot:
	goreleaser release --snapshot --clean

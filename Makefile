.PHONY: test lint


default: test lint

test:
	go test ./...

lint:
	go tool revive -config revive.toml --formatter friendly --exclude *_test.go ./...

pretty:
	go tool gofumpt -l -w .
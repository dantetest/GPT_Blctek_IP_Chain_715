.PHONY: fmt test vet check run-api run-worker run-agent infra-up infra-down

fmt:
	gofmt -w $$(find apps packages -name '*.go' -type f)

test:
	go test ./apps/api/... ./apps/worker/... ./apps/data-agent/... ./packages/manifest-spec/...

vet:
	go vet ./apps/api/... ./apps/worker/... ./apps/data-agent/... ./packages/manifest-spec/...

check: fmt test vet

run-api:
	go run ./apps/api/cmd/api

run-worker:
	go run ./apps/worker/cmd/worker

run-agent:
	go run ./apps/data-agent/cmd/agent --help

infra-up:
	docker compose up -d mysql redis minio

infra-down:
	docker compose down

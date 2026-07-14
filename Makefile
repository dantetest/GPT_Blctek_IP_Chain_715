.PHONY: fmt test vet check run-api run-worker infra-up infra-down

fmt:
	gofmt -w $$(find apps -name '*.go' -type f)

test:
	go test ./apps/api/... ./apps/worker/...

vet:
	go vet ./apps/api/... ./apps/worker/...

check: fmt test vet

run-api:
	go run ./apps/api/cmd/api

run-worker:
	go run ./apps/worker/cmd/worker

infra-up:
	docker compose up -d mysql redis minio

infra-down:
	docker compose down

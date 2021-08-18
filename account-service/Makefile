PHONY: build

build:
	@go build -ldflags="-s -w" -o account ./main.go
	@echo "Binary generated."

run:
	@go run ./main.go

tests:
	@go test -v ./test/...

coverage:
	@go test -v -cover -coverprofile=r.out -coverpkg ./internal/... ./test/...
	@go tool cover -html=r.out -o report.html
	@rm -f r.out

docker:
	@docker build -t account -f deployments/Dockerfile .
	@echo "Docker image generated."

prune:
	@docker image prune --filter label=stage=builder
	@echo "Removed build intermediate images."

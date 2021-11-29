.PHONY: dockerrun run test doc lint

TAG=neemiasjnr/echod

dockerrun:
	@docker build -t ${TAG} -f Dockerfile-app . && docker run -p 3000:3000 ${TAG}

# Recommended: make run
run:
	@docker-compose up

test:
	@go test -v ./internal/server

doc:
	@echo "Access: http://localhost:6060/pkg/$(shell go list -m)/internal/server" && godoc -http=:6060 >/dev/null 2>&1

lint:
	@golangci-lint run && go vet && staticcheck

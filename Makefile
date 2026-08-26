APP=mpc-coordinator
.PHONY: all fmt test race vet build run smoke
all: fmt test vet build
fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')
test:
	go test ./...
race:
	go test -race ./...
vet:
	go vet ./...
build:
	go build -trimpath ./...
run:
	MPC_STATE_FILE=./data/state.json go run ./cmd/mpc-coordinator
smoke:
	./scripts/smoke.sh


GO_VERSION:=1.23


.PHONY: build
	

build: 
	@go build ./... 

build-optimized:
	@go build  -ldflags="-s -w" -o $(GOPATH)/bin/dolus-optimized cmd/dolus/main.go

install: build
	@go install ./...

run: 
	@go run cmd/dolus-test/main.go

run-cli:
	@go run cmd/dolus/main.go

run-optimized: build-optimized
	@dolus-optimized

debug: $(GOPATH)/bin/dlv
	@dlv debug cmd/dolus-test/main.go

size:
	@du -h $(GOPATH)/bin/dolus

size-optimized:
	@du -h $(GOPATH)/bin/dolus-optimized


# --- DOCKER ---

docker-build:
	docker build -t dolus -t dolusmockserver/dolus:latest -f build/package/Dockerfile .

docker-push:
	docker push dolusmockserver/dolus:latest


gen: gen-go-server-client gen-cue-expectations gen-mocks

gen-go-server-client: $(GOPATH)/bin/oapi-codegen
	oapi-codegen --package=api -generate=server,types,spec,client api/dolus.yaml > internal/api/api.gen.go
	

gen-cue-expectations:
	cd cue-expectations && cue get go ../pkg/expectation/ -e \
	ExpectationError,ExpectationFieldError,Route


gen-mocks:
	mockery

test:
	go test ./...

test-verbose:
	go test -v ./...

test-nice:
	go test ./...  -json | tparse -all


### TOOLS ###
$(GOPATH)/bin/dlv:
	go install github.com/go-delve/delve/cmd/dlv  

$(GOPATH)/bin/oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest    


update-dolus-expectations:
	./install-local.sh

lint-api-spec:
	spectral lint dolus.yaml


### DAGGER ###

dagger-test:
	dagger call -m golang test --source . --go-version=$GO_VERSION stdout



### ACT ###
act:
	act -P ubuntu-latest=catthehacker/ubuntu:act-latest





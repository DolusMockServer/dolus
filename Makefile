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


gen: gen-go-server-client gen-cue-expectations

gen-go-server-client: $(GOPATH)/bin/oapi-codegen
	oapi-codegen --package=api -generate=server,types,spec,client api/dolus.yaml > internal/api/api.gen.go
	

gen-cue-expectations:
	go run cmd/cue2gostruct/main.go cue-expectations/core/core.cue pkg/expectation/cue/expectation.gen.go \
	cue
	go fmt pkg/expectation/cue/expectation.gen.go


### TOOLS ###
$(GOPATH)/bin/dlv:
	go install github.com/go-delve/delve/cmd/dlv  

$(GOPATH)/bin/oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest    



update-dolus-expectations:
	./install-local.sh

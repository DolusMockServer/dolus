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


gen: $(GOPATH)/bin/oapi-codegen
	oapi-codegen --package=server -generate=server,types,spec api/dolus.yaml > internal/server/server.gen.go


gen-go-sdk-client: $(GOPATH)/bin/oapi-codegen
	oapi-codegen --package=client -generate=client,types api/dolus.yaml > internal/client/client.gen.go 
	

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

.PHONY: deps
deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

.PHONY: all build
all build:
	@go build -o build/actor-stdout cmd/actor-stdout/main.go 
	@go build -o build/leappctl cmd/leappctl/main.go 
	@go build -o build/leapp-daemon cmd/leapp-daemon/main.go 

.PHONY: test
test:
	@go test

.PHONY: clean
clean:
	@rm -rf build/

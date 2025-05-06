.PHONY: default run build release mkdirs clean 
PLATFORM=$(shell sh -c "go version | awk '{print \$$4}' | tr '/' '-'")
CMD=hs
PACKAGE=github.com/hurricanerix/http-helper
BUILD_DIR=bin
HOME_DIR=`echo $(HOME)`

SRC=$(shell find . -type f -regex ".*\.go")
BASE64_SOURCE_DIFF=$(shell git --no-pager diff | base64 -w0)
VERSION_TAG=$(shell git describe --exact-match --tags || echo "")

.PHONY: default
default: build

.PHONY: run
run: build
	./bin/$(PLATFORM)/hs -d testdata

.PHONY: setup
setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: lint
lint:
	go vet ./...
	$(GOPATH)/bin/staticcheck -checks all ./...

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	@echo "Not implemented yet"

.PHONY: fuzz
fuzz:
	@echo "Not implemented yet"

.PHONY: build
build: $(BUILD_DIR) $(BUILD_DIR)/$(PLATFORM)/$(CMD)

.PHONY: build-all
build-all: $(BUILD_DIR) $(BUILD_DIR)/darwin-arm64/$(CMD) $(BUILD_DIR)/darwin-amd64/$(CMD) $(BUILD_DIR)/windows-amd64/$(CMD).exe $(BUILD_DIR)/linux-amd64/$(CMD)


$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/darwin-arm64:
	mkdir -p $(BUILD_DIR)/darwin-arm64

$(BUILD_DIR)/darwin-amd64:
	mkdir -p $(BUILD_DIR)/darwin-amd64

$(BUILD_DIR)/windows-amd64:
	mkdir -p $(BUILD_DIR)/windows-amd64

$(BUILD_DIR)/linux-amd64:
	mkdir -p $(BUILD_DIR)/linux-amd64

$(BUILD_DIR)/linux-amd64/$(CMD): $(BUILD_DIR)/linux-amd64 $(SRC)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X $(PACKAGE)/build.base64SourceDiff=$(BASE64_SOURCE_DIFF) -X $(PACKAGE)/build.version=$(VERSION_TAG)" -o ./$@ $(PACKAGE)/cmd/$(CMD)

$(BUILD_DIR)/darwin-arm64/$(CMD): $(BUILD_DIR)/darwin-arm64 $(SRC)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "-X $(PACKAGE)/build.base64SourceDiff=$(BASE64_SOURCE_DIFF) -X $(PACKAGE)/build.version=$(VERSION_TAG)"" -o ./$@ $(PACKAGE)/cmd/$(CMD)

$(BUILD_DIR)/darwin-amd64/$(CMD): $(BUILD_DIR)/darwin-amd64 $(SRC)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X $(PACKAGE)/build.base64SourceDiff=$(BASE64_SOURCE_DIFF) -X $(PACKAGE)/build.version=$(VERSION_TAG)"" -o ./$@ $(PACKAGE)/cmd/$(CMD)

$(BUILD_DIR)/windows-amd64/$(CMD).exe: $(BUILD_DIR)/windows-amd64 $(SRC)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X $(PACKAGE)/build.base64SourceDiff=$(BASE64_SOURCE_DIFF) -X $(PACKAGE)/build.version=$(VERSION_TAG)"" -o ./$@ -buildmode=exe $(PACKAGE)/cmd/$(CMD)

.PHONY: release
release: build
	@tar czf "$(BUILD_DIR)/$(CMD)-linux-amd64.tar.gz" --directory="$(BUILD_DIR)/linux-amd64" "$(CMD)"
	@cd $(BUILD_DIR); shasum -a 256  "$(CMD)-linux-amd64.tar.gz"
	@zip -q -r -j "$(BUILD_DIR)/$(CMD)-windows-amd64.zip" "$(BUILD_DIR)/windows-amd64/$(CMD).exe"
	@cd $(BUILD_DIR); shasum -a 256 "$(CMD)-windows-amd64.zip"
	@tar czf "$(BUILD_DIR)/$(CMD)-darwin-amd64.tar.gz" --directory="$(BUILD_DIR)/darwin-amd64" "$(CMD)"
	@cd $(BUILD_DIR); shasum -a 256 "$(CMD)-darwin-amd64.tar.gz"
	@tar czf "$(BUILD_DIR)/$(CMD)-darwin-arm64.tar.gz" --directory="$(BUILD_DIR)/darwin-arm64" "$(CMD)"
	@cd $(BUILD_DIR); shasum -a 256 "$(CMD)-darwin-arm64.tar.gz"

.PHONY: clean
clean: 
	rm -rf $(BUILD_DIR)/linux-amd64
	rm -rf $(BUILD_DIR)/windows-amd64
	rm -rf $(BUILD_DIR)/darwin-amd64
	rm -rf $(BUILD_DIR)/darwin-arm64
	rm -f $(BUILD_DIR)/$(CMD)-linux-amd64.tar.gz
	rm -f $(BUILD_DIR)/$(CMD)-windows-amd64.zip
	rm -f $(BUILD_DIR)/$(CMD)-darwin-amd64.tar.gz
	rm -f $(BUILD_DIR)/$(CMD)-darwin-arm64.tar.gz

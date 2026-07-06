# likhis - Universal API Route Mapper
# Makefile for clean building across platforms

BINARY_NAME    = likhis
BUILD_DIR      = build
OUT_DIR        = out
RELEASE_DIR    = release
PLUGIN_DIR     = plugins
LINK_DIR       = C:\Bin\tools
GOFLAGS        = -ldflags="-s -w"
GOTEST_FLAGS   = -v

ifneq ($(MSYSTEM),)
    # MSYS2/Git Bash - Unix commands available
    BINARY      = $(BINARY_NAME).exe
    RM          = rm -rf
    RM_FILE     = rm -f
    MKDIR_P     = mkdir -p
    GO_BUILD    = go build $(GOFLAGS) -o
    TEST_RUN    = go test ./tests $(GOTEST_FLAGS)
    TEST_COV    = go test ./tests -coverprofile=coverage.out
else ifeq ($(OS),Windows_NT)
    # Native Windows cmd
    BINARY      = $(BINARY_NAME).exe
    SHELL       = cmd.exe
    RM          = rmdir /s /q 2>nul
    RM_FILE     = del /f /q 2>nul
    MKDIR_P     = mkdir 2>nul
    GO_BUILD    = go build $(GOFLAGS) -o
    TEST_RUN    = go test ./tests $(GOTEST_FLAGS)
    TEST_COV    = go test ./tests -coverprofile=coverage.out
else
    BINARY      = $(BINARY_NAME)
    RM          = rm -rf
    RM_FILE     = rm -f
    MKDIR_P     = mkdir -p
    GO_BUILD    = go build $(GOFLAGS) -o
    TEST_RUN    = go test ./tests $(GOTEST_FLAGS)
    TEST_COV    = go test ./tests -coverprofile=coverage.out
endif

.PHONY: all build clean test test-integration coverage lint run release link help

all: clean build

build: $(BUILD_DIR)/$(BINARY)

$(BUILD_DIR)/$(BINARY):
	@echo Building $(BINARY)...
	-$(MKDIR_P) $(BUILD_DIR)
	$(GO_BUILD) $(BUILD_DIR)/$(BINARY) main.go
	@echo ✓ Built $(BUILD_DIR)/$(BINARY)

clean:
	@echo Cleaning build artifacts...
	-$(RM) $(BUILD_DIR)
	-$(RM) $(OUT_DIR)
	-$(RM) $(RELEASE_DIR)
	-$(RM_FILE) coverage.out
	-$(RM_FILE) *.test
	@echo ✓ Cleaned

test:
	@echo Running unit tests...
	$(TEST_RUN)
	@echo ✓ Unit tests passed

test-integration: build
	@echo Running integration tests against example projects...
	$(MKDIR_P) $(OUT_DIR)
	$(MKDIR_P) $(OUT_DIR)/express
	$(MKDIR_P) $(OUT_DIR)/flask
	$(MKDIR_P) $(OUT_DIR)/django
	$(MKDIR_P) $(OUT_DIR)/laravel
	$(MKDIR_P) $(OUT_DIR)/spring
	$(BUILD_DIR)/$(BINARY) -p exp/express -o postman -F express -O $(OUT_DIR)/express && echo [PASSED] Express
	$(BUILD_DIR)/$(BINARY) -p exp/flask -o postman -F flask -O $(OUT_DIR)/flask && echo [PASSED] Flask
	$(BUILD_DIR)/$(BINARY) -p exp/django -o postman -F django -O $(OUT_DIR)/django && echo [PASSED] Django
	$(BUILD_DIR)/$(BINARY) -p exp/laravel -o postman -F laravel -O $(OUT_DIR)/laravel && echo [PASSED] Laravel
	$(BUILD_DIR)/$(BINARY) -p exp/spring -o postman -F spring -O $(OUT_DIR)/spring && echo [PASSED] Spring
	$(BUILD_DIR)/$(BINARY) -p exp/express -o postman -F auto -O $(OUT_DIR)/express && echo [PASSED] Auto-detect
	$(BUILD_DIR)/$(BINARY) -p exp/express -o postman -F express --full -O $(OUT_DIR)/express && echo [PASSED] Full export
	$(BUILD_DIR)/$(BINARY) -p exp/express -o insomnia -F express -O $(OUT_DIR)/express && echo [PASSED] Insomnia
	$(BUILD_DIR)/$(BINARY) -p exp/express -o httpie -F express -O $(OUT_DIR)/express && echo [PASSED] HTTPie
	$(BUILD_DIR)/$(BINARY) -p exp/express -o curl -F express -O $(OUT_DIR)/express && echo [PASSED] cURL
	@echo ✓ Integration tests passed

coverage:
	@echo Running tests with coverage...
	$(TEST_COV)
	go tool cover -func=coverage.out
	@echo ✓ Coverage report generated

lint:
	@echo Running linter...
	golangci-lint run --timeout=5m --skip-files-use-default-excludes
	@echo ✓ Lint passed

run: build
	@echo Running likhis on current directory...
	$(BUILD_DIR)/$(BINARY) -p . -o postman

release: clean
	@echo Building release binaries...
	$(MKDIR_P) $(RELEASE_DIR)
ifneq ($(MSYSTEM),)
	GOOS=windows GOARCH=amd64 $(GO_BUILD) $(RELEASE_DIR)/likhis-windows-amd64.exe main.go
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) $(RELEASE_DIR)/likhis-linux-amd64   main.go
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) $(RELEASE_DIR)/likhis-darwin-amd64  main.go
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) $(RELEASE_DIR)/likhis-darwin-arm64  main.go
else ifeq ($(OS),Windows_NT)
	set GOOS=windows&& set GOARCH=amd64&& $(GO_BUILD) $(RELEASE_DIR)/likhis-windows-amd64.exe main.go
	set GOOS=linux&&   set GOARCH=amd64&& $(GO_BUILD) $(RELEASE_DIR)/likhis-linux-amd64   main.go
	set GOOS=darwin&&  set GOARCH=amd64&& $(GO_BUILD) $(RELEASE_DIR)/likhis-darwin-amd64  main.go
	set GOOS=darwin&&  set GOARCH=arm64&& $(GO_BUILD) $(RELEASE_DIR)/likhis-darwin-arm64  main.go
else
	GOOS=windows GOARCH=amd64 $(GO_BUILD) $(RELEASE_DIR)/likhis-windows-amd64.exe main.go
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) $(RELEASE_DIR)/likhis-linux-amd64   main.go
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) $(RELEASE_DIR)/likhis-darwin-amd64  main.go
	GOOS=darwin  GOARCH=arm64 $(GO_BUILD) $(RELEASE_DIR)/likhis-darwin-arm64  main.go
endif
	@echo ✓ Release binaries built in $(RELEASE_DIR)/

link: build
ifneq ($(MSYSTEM),)
	@echo Creating symbolic link at $(LINK_DIR)...
	mkdir -p $(LINK_DIR)
	ln -sf $(CURDIR)/$(BUILD_DIR)/$(BINARY) $(LINK_DIR)/$(BINARY)
	@echo ✓ Symbolic link created: $(LINK_DIR)/$(BINARY)
else ifeq ($(OS),Windows_NT)
	@echo Creating symbolic link at $(LINK_DIR)...
	if not exist $(LINK_DIR) mkdir $(LINK_DIR)
	mklink "$(LINK_DIR)\$(BINARY)" "%CD%\$(BUILD_DIR)\$(BINARY)"
	@echo ✓ Symbolic link created: $(LINK_DIR)\$(BINARY)
else
	@echo Creating symbolic link...
	ln -sf $(CURDIR)/$(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)
	@echo ✓ Symbolic link created: /usr/local/bin/$(BINARY)
endif

help:
	@echo Usage: make [target]
	@echo.
	@echo Targets:
	@echo   all              Clean and build (default)
	@echo   build            Build the binary
	@echo   clean            Remove build artifacts
	@echo   test             Run unit tests
	@echo   test-integration Run integration tests against example projects
	@echo   coverage         Run tests with coverage report
	@echo   lint             Run golangci-lint
	@echo   run              Build and run on current directory
	@echo   release          Build release binaries for all platforms
	@echo   link             Create symbolic link to PATH
	@echo   help             Show this help message

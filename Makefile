# Go related variables
GO = go
GOFMT = gofmt
GOSTATICCHECK = go run honnef.co/go/tools/cmd/staticcheck@v0.4.3
GOMOCKERY = $(GO) run github.com/vektra/mockery/v2@v2.32.4

M = $(shell printf "\033[34;1m▶▶▶\033[0m")

# Test coverage files
TEST_COVERAGE_PROFILE_OUTPUT = ".local/coverage.out"
TEST_REPORT_OUTPUT = ".local/test_report.ndjson"

all: verify

.PHONY: verify
verify:
	$(info $(M) Verifying...) @
	chmod +x ./verify.sh
	./verify.sh

# -----------------------------------------------------------------------------
# Build

.PHONY: build
build:
	$(info $(M) Building...) @
	$(GO) build -o bin/ ./...

.PHONY: cross-build
cross-build:
	$(info $(M) Building...) @
	./dist.sh

.PHONY: install
install:
	$(info $(M) Installing...) @
	$(GO) install ./...

# -----------------------------------------------------------------------------
# Test

.PHONY: test
test:
	$(info $(M) Running tests...) @
	$(GO) test ./...

.PHONY: coverage
coverage:
	$(info $(M) Generating coverage information...)
	$(eval TEST_COVERAGE_PROFILE_OUTPUT_DIRNAME=$(shell dirname $(TEST_COVERAGE_PROFILE_OUTPUT)))
	$(eval TEST_REPORT_OUTPUT_DIRNAME=$(shell dirname $(TEST_REPORT_OUTPUT)))
	mkdir -p $(TEST_COVERAGE_PROFILE_OUTPUT_DIRNAME) $(TEST_REPORT_OUTPUT_DIRNAME)
	$(GO) test ./... -coverpkg=./... -coverprofile=$(TEST_COVERAGE_PROFILE_OUTPUT) -json > $(TEST_REPORT_OUTPUT)

.PHONY: generate-mocks
generate-mocks:
	$(info $(M) Running mockery...)
	$(GOMOCKERY)

# ----------------------------------------------------------
# Dependencies

## tidy: Executes a go mod tidy
.PHONY: tidy
tidy:
	$(info $(M) Running go mod tidy...) @
	$(GO) mod tidy

# ----------------------------------------------------------
# Format and lint

## checkfmt: Check format validation
.PHONY: checkfmt
checkfmt:
	$(info $(M) Running gofmt checking code style...)
	@fmtRes=$$($(GOFMT) -d .); \
	if [ -n "$${fmtRes}" ]; then \
		echo "gofmt checking failed!"; echo "$${fmtRes}"; echo; \
		echo "Please ensure you are using $$($(GO) version) for formatting code."; \
		exit 1; \
	fi

## fmt: Formats the code
.PHONY: fmt
fmt:
	$(info $(M) Running go fmt...) @
	$(GOFMT) -l -w .


## lint: Runs go linter
.PHONY: lint
lint:
	$(info $(M) Running lints...)
	$(GOSTATICCHECK) ./...

## vet: Run go vet
.PHONY: vet
vet:
	$(info $(M) Running go vet...)
	$(GO) vet ./...

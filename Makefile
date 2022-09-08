.PHONY: all build clean update static static-test-server test test-server version-tag build-docker

APP_NAME = go-microservices
SWAGGER_FILE = openapi-go-microservices.json
GO_MICROSERVICES_TEST_PORT ?= 3000
INSTALL_BASE ?= /usr/local

# ==========================================================================================================================

GO_TEST_PACKAGES = service
# make phony targets out of the go package/subdirectory list
GO_TEST_PACKAGE_TARGETS = $(addprefix gotest_,$(GO_TEST_PACKAGES))

# Find all Go files for this service
GO_FILES := $(shell find . -type f -name '*.go')

# ==========================================================================================================================

# Filter out and only accept tags that start with v. Anything else gets "latest".
VERSION := $(if $(filter v%,$(CI_COMMIT_TAG)),$(CI_COMMIT_TAG:v%=%),latest)
# git commit hash
CI_COMMIT_SHA := $(shell git describe --dirty --long 2>/dev/null || echo 'No_CI_System')
# RPM Release Version
RELEASE = 1

# ==========================================================================================================================

all: build

test: test-server

static: static-test-server

build: build-server

push: push-docker

swagger: swagger

# ==========================================================================================================================

## UNIT TESTS
test-server: $(GO_TEST_PACKAGE_TARGETS)
	go test -v

$(GO_TEST_PACKAGE_TARGETS):
	go test -v "./$(@:gotest_%=%)"

static-test-server:
	golangci-lint run

# ==========================================================================================================================

## PROJECT BUILD
build-server:
	go build -ldflags "-X main.buildVersion=$(VERSION) -X main.gitCommit=$(CI_COMMIT_SHA)" -o $(APP_NAME) main.go

# ==========================================================================================================================

## DOCKER BUILD

build-docker: build build-tar
	docker build -t $(APP_NAME):$(VERSION) .

# ==========================================================================================================================

## REMOVE PREVIOUS BUILD
clean:
	rm -f $(APP_NAME)

# ==========================================================================================================================

## SWAGGER API SPEC GENERATION
# Setting SWAGGER_GENERATE_EXTENSION to false in the environment skips generation of any x-go fields
swagger:
	SWAGGER_GENERATE_EXTENSION="false" swagger generate spec -o $(SWAGGER_FILE) --scan-models

GO ?= go

WORK_DIR   := $(shell pwd)

KHULNASOFT_SDK_TEST_URL ?= http://localhost:3000
KHULNASOFT_SDK_TEST_USERNAME ?= test01
KHULNASOFT_SDK_TEST_PASSWORD ?= test01

PACKAGE := github.com/khulnasoft/go-sdk

GOFUMPT_PACKAGE ?= mvdan.cc/gofumpt@v0.7.0
GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
KHULNASOFT_VET_PACKAGE ?= github.com/khulnasoft/go-sdk/khulnasoft-vet@latest

KHULNASOFT_VERSION := 1.21.10
KHULNASOFT_DL := https://dl.khulnasoft.com/khulnasoft/$(KHULNASOFT_VERSION)/khulnasoft-$(KHULNASOFT_VERSION)-
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
  KHULNASOFT_DL := $(KHULNASOFT_DL)linux-

  UNAME_P := $(shell uname -p)
  ifeq ($(UNAME_P),unknown)
   KHULNASOFT_DL := $(KHULNASOFT_DL)amd64
  endif
  ifeq ($(UNAME_P),x86_64)
   KHULNASOFT_DL := $(KHULNASOFT_DL)amd64
  endif
  ifneq ($(filter %86,$(UNAME_P)),)
   KHULNASOFT_DL := $(KHULNASOFT_DL)386
  endif
  ifneq ($(filter arm%,$(UNAME_P)),)
    KHULNASOFT_DL := $(KHULNASOFT_DL)arm-5
  endif
endif
ifeq ($(UNAME_S),Darwin)
  KHULNASOFT_DL := $(KHULNASOFT_DL)darwin-10.12-amd64
endif

.PHONY: all
all: clean test build

.PHONY: help
help:
	@echo "Make Routines:"
	@echo " - \"\"              run \"make clean test build\""
	@echo " - build             build sdk"
	@echo " - clean             clean"
	@echo " - fmt               format the code"
	@echo " - lint              run golint"
	@echo " - vet               examines Go source code and reports"
	@echo " - test              run unit tests (need a running khulnasoft)"
	@echo " - test-instance     start a khulnasoft instance for test"


.PHONY: clean
clean:
	rm -r -f test
	cd khulnasoft && $(GO) clean -i ./...

.PHONY: fmt
fmt:
	find . -name "*.go" -type f | xargs gofmt -s -w; \
	$(GO) run $(GOFUMPT_PACKAGE) -extra -w ./khulnasoft

.PHONY: vet
vet: build
	 cd khulnasoft-vet $(GO) vet ./...
	 cd khulnasoft-vet $(GO) vet -vettool=khulnasoft-vet ./...

.PHONY: ci-lint
ci-lint:
	@cd khulnasoft/; echo -n "gofumpt ...";\
	diff=$$($(GO) run $(GOFUMPT_PACKAGE) -extra -l .); \
	if [ -n "$$diff" ]; then \
		echo; echo "Not gofumpt-ed"; \
		exit 1; \
	fi; echo " done"; echo -n "golangci-lint ...";\
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run --timeout 5m; \
	if [ $$? -eq 1 ]; then \
		echo; echo "Doesn't pass golangci-lint"; \
		exit 1; \
	fi; echo " done"; \
	cd -; \

.PHONY: test
test:
	@export KHULNASOFT_SDK_TEST_URL=${KHULNASOFT_SDK_TEST_URL}; export KHULNASOFT_SDK_TEST_USERNAME=${KHULNASOFT_SDK_TEST_USERNAME}; export KHULNASOFT_SDK_TEST_PASSWORD=${KHULNASOFT_SDK_TEST_PASSWORD}; \
	if [ -z "$(shell curl --noproxy "*" "${KHULNASOFT_SDK_TEST_URL}/api/v1/version" 2> /dev/null)" ]; then \echo "No test-instance detected!"; exit 1; else \
	    cd khulnasoft && $(GO) test -race -cover -coverprofile coverage.out; \
	fi

.PHONY: test-instance
test-instance:
	rm -f -r ${WORK_DIR}/test 2> /dev/null; \
	mkdir -p ${WORK_DIR}/test/conf/ ${WORK_DIR}/test/data/
	wget ${KHULNASOFT_DL} -O ${WORK_DIR}/test/khulnasoft-main; \
	chmod +x ${WORK_DIR}/test/khulnasoft-main; \
	echo "[security]" > ${WORK_DIR}/test/conf/app.ini; \
	echo "INTERNAL_TOKEN = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NTg4MzY4ODB9.LoKQyK5TN_0kMJFVHWUW0uDAyoGjDP6Mkup4ps2VJN4" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "INSTALL_LOCK   = true" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "SECRET_KEY     = 2crAW4UANgvLipDS6U5obRcFosjSJHQANll6MNfX7P0G3se3fKcCwwK3szPyGcbo" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "PASSWORD_COMPLEXITY = off" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[database]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "DB_TYPE = sqlite3" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[repository]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "ROOT = ${WORK_DIR}/test/data/" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[server]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "ROOT_URL = ${KHULNASOFT_SDK_TEST_URL}" >> ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/khulnasoft-main migrate -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/khulnasoft-main admin user create --username=${KHULNASOFT_SDK_TEST_USERNAME} --password=${KHULNASOFT_SDK_TEST_PASSWORD} --email=test01@khulnasoft.io --admin=true --must-change-password=false --access-token -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/khulnasoft-main web -c ${WORK_DIR}/test/conf/app.ini

.PHONY: bench
bench:
	cd khulnasoft && $(GO) test -run=XXXXXX -benchtime=10s -bench=. || exit 1

.PHONY: build
build:
	cd khulnasoft && $(GO) build


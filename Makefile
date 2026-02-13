# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

ENSURE_GARDENER_MOD         := $(shell go get github.com/gardener/gardener@$$(go list -m -f "{{.Version}}" github.com/gardener/gardener))
GARDENER_HACK_DIR           := $(shell go list -m -f "{{.Dir}}" github.com/gardener/gardener)/hack
NAME                        := garden-shoot-trust-configurator
IMAGE                       := europe-docker.pkg.dev/gardener-project/public/gardener/$(NAME)
REPO_ROOT                   := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
HACK_DIR                    := $(REPO_ROOT)/hack
VERSION                     := $(shell cat "$(REPO_ROOT)/VERSION")
GOARCH                      ?= $(shell go env GOARCH)
EFFECTIVE_VERSION           := $(VERSION)-$(shell git rev-parse HEAD)
LD_FLAGS                    := "-w $(shell bash $(GARDENER_HACK_DIR)/get-build-ld-flags.sh k8s.io/component-base $(REPO_ROOT)/VERSION $(NAME))"

ifneq ($(strip $(shell git status --porcelain 2>/dev/null)),)
	EFFECTIVE_VERSION := $(EFFECTIVE_VERSION)-dirty
endif

TOOLS_DIR := $(REPO_ROOT)/hack/tools
include $(GARDENER_HACK_DIR)/tools.mk

.PHONY: start
start:
	go run -ldflags $(LD_FLAGS) ./cmd/garden-shoot-trust-configurator/main.go --config=./example/00-config.yaml

#################################################################
# Rules related to binary build, Docker image build and release #
#################################################################

.PHONY: install
install:
	@LD_FLAGS=$(LD_FLAGS) EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) \
		bash $(GARDENER_HACK_DIR)/install.sh ./...

.PHONY: docker-images
docker-images:
	@docker build --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) --build-arg TARGETARCH=$(GOARCH) -t $(IMAGE):$(EFFECTIVE_VERSION) -t $(IMAGE):latest -f Dockerfile --target $(NAME) . --memory 6g

#####################################################################
# Rules for verification, formatting, linting, testing and cleaning #
#####################################################################

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: clean
clean:
	@rm -f $(REPO_ROOT)/gosec-report.sarif
	@bash $(GARDENER_HACK_DIR)/clean.sh ./cmd/... ./internal/... ./pkg/...

.PHONY: check-generate
check-generate:
	@bash $(GARDENER_HACK_DIR)/check-generate.sh $(REPO_ROOT)

.PHONY: check
check: $(GOIMPORTS) $(GOLANGCI_LINT) $(HELM) $(YQ) $(TYPOS) 
	go vet ./...
	@REPO_ROOT=$(REPO_ROOT) bash $(GARDENER_HACK_DIR)/check.sh --golangci-lint-config=./.golangci.yaml ./cmd/... ./internal/... ./pkg/...
	@bash $(GARDENER_HACK_DIR)/check-typos.sh
	@bash $(GARDENER_HACK_DIR)/check-file-names.sh
	@bash $(GARDENER_HACK_DIR)/check-charts.sh ./charts

tools-for-generate: $(CONTROLLER_GEN) $(YQ) $(MOCKGEN) $(HELM) $(GEN_CRD_API_REFERENCE_DOCS)
	@go mod download

.PHONY: generate
generate: tools-for-generate
	@REPO_ROOT=$(REPO_ROOT) GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) bash $(GARDENER_HACK_DIR)/generate-sequential.sh ./cmd/... ./internal/... ./pkg/...
	@REPO_ROOT=$(REPO_ROOT) GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) $(REPO_ROOT)/hack/update-codegen.sh
	$(MAKE) format

.PHONY: format
format: $(GOIMPORTS) $(GOIMPORTSREVISER)
	@bash $(GARDENER_HACK_DIR)/format.sh ./cmd ./internal ./pkg

.PHONY: sast
sast: $(GOSEC)
	@bash $(GARDENER_HACK_DIR)/sast.sh --exclude-dirs dev

.PHONY: sast-report
sast-report: $(GOSEC)
	@bash $(GARDENER_HACK_DIR)/sast.sh --gosec-report true --exclude-dirs dev

.PHONY: test
test: $(REPORT_COLLECTOR)
	@bash $(GARDENER_HACK_DIR)/test.sh ./cmd/... ./internal/... ./pkg/...	

.PHONY: test-cov
test-cov:
	@bash $(GARDENER_HACK_DIR)/test-cover.sh ./cmd/... ./internal/... ./pkg/...

.PHONY: test-clean
test-clean:
	@bash $(GARDENER_HACK_DIR)/test-cover-clean.sh

.PHONY: verify
verify: check format test sast

.PHONY: verify-extended
verify-extended: check-generate check format test test-cov test-clean sast-report

##############################################################
# Rules related to kind and skaffold based local development #
##############################################################

export SKAFFOLD_BUILD_CONCURRENCY = 0
server-up server-down: export SKAFFOLD_DEFAULT_REPO = garden.local.gardener.cloud:5001
server-up server-down: export SKAFFOLD_PUSH = true
# use static label for skaffold to prevent rolling all gardener components on every `skaffold` invocation
server-up server-down: export SKAFFOLD_LABEL = skaffold.dev/run-id=server-local

server-up: $(SKAFFOLD) $(KIND) $(HELM) $(KUBECTL)
	@LD_FLAGS=$(LD_FLAGS) GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) $(SKAFFOLD) run

server-down: $(SKAFFOLD) $(HELM) $(KUBECTL)
	$(SKAFFOLD) delete

## CI E2E Tests
ci-e2e-kind:
	./hack/ci-e2e-kind.sh

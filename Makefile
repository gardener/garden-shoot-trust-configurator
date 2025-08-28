# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

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
	go run -ldflags $(LD_FLAGS) ./cmd/garden-shoot-trust-configurator/main.go 

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
	@bash $(GARDENER_HACK_DIR)/clean.sh ./cmd/... ./internal/...

.PHONY: check-generate
check-generate:
	@bash $(GARDENER_HACK_DIR)/check-generate.sh $(REPO_ROOT)

.PHONY: check
check: $(GOIMPORTS) $(GOLANGCI_LINT) $(HELM) $(YQ)
	@bash $(GARDENER_HACK_DIR)/check.sh --golangci-lint-config=./.golangci.yaml ./cmd/... ./internal/...
#TODO(theoddora): re-enable when adding skaffold based local dev setup + charts
# 	@bash $(GARDENER_HACK_DIR)/check-charts.sh ./charts
# 	@GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) hack/check-skaffold-deps.sh

# .PHONY: update-skaffold-deps
# update-skaffold-deps: $(YQ)
# 	@GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) hack/check-skaffold-deps.sh update

.PHONY: generate
generate: $(CONTROLLER_GEN) $(YQ) $(VGOPATH) $(MOCKGEN) $(HELM)
#TODO(theoddora): add charts directory when charts are introduced to the repository
	@VGOPATH=$(VGOPATH) REPO_ROOT=$(REPO_ROOT) GARDENER_HACK_DIR=$(GARDENER_HACK_DIR) bash $(GARDENER_HACK_DIR)/generate-sequential.sh ./cmd/... ./internal/...
	$(MAKE) format

.PHONY: format
format: $(GOIMPORTS) $(GOIMPORTSREVISER)
	@bash $(GARDENER_HACK_DIR)/format.sh ./cmd ./internal

.PHONY: sast
sast: $(GOSEC)
	@bash $(GARDENER_HACK_DIR)/sast.sh

.PHONY: sast-report
sast-report: $(GOSEC)
	@bash $(GARDENER_HACK_DIR)/sast.sh --gosec-report true

.PHONY: test
test: $(REPORT_COLLECTOR)
	@bash $(GARDENER_HACK_DIR)/test.sh ./cmd/... ./internal/...

.PHONY: test-cov
test-cov:
	@bash $(GARDENER_HACK_DIR)/test-cover.sh ./cmd/... ./internal/...g

.PHONY: test-clean
test-clean:
	@bash $(GARDENER_HACK_DIR)/test-cover-clean.sh

.PHONY: verify
verify: check format test sast

.PHONY: verify-extended
# TODO(theoddora): add test command after including unit/integration tests
verify-extended: check-generate check format sast-report

# TODO(theoddora): re-enable when adding skaffold based local dev setup
# # use static label for skaffold to prevent rolling all gardener components on every `skaffold` invocation
# server-up server-down: export SKAFFOLD_LABEL = skaffold.dev/run-id=server-local

# server-up: $(SKAFFOLD) $(KIND) $(HELM)
# 	@LD_FLAGS=$(LD_FLAGS) $(SKAFFOLD) run

# server-dev: $(SKAFFOLD) $(HELM)
# 	$(SKAFFOLD) dev --cleanup=false --trigger=manual

# server-down: $(SKAFFOLD) $(HELM)
# 	$(SKAFFOLD) delete

ci-e2e-kind:
	./hack/ci-e2e-kind.sh

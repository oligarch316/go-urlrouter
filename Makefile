#!/usr/bin/make -f

# Disable all default suffixes
.SUFFIXES:

cmd_go := $(shell command -v go1.18beta1 || echo "go")


# ----- Tooling
.PHONY: tidy

tidy: ## Tidy go modules
	$(info Tidying)
	@$(cmd_go) mod tidy


# ----- Test
test_flags :=
test_coverprofile_target := coverage.out

.PHONY: test test.v test.cover coverage test.clean

test: ## Run tests
	$(info Testing)
	@$(cmd_go) test $(test_flags) ./...

test.v: test_flags += -v
test.v: test ## Run tests with verbose output

test.cover: test_flags += -coverprofile=$(test_coverprofile_target)
test.cover: test ## Generate test coverage report

coverage: test.cover ## Generate and open test coverage report in browser
	@$(cmd_go) tool cover -html $(test_coverprofile_target)

test.clean: ## Clean test artifacts
	$(info Cleaning test artifacts)
	@rm $(test_coverprofile_target) 2> /dev/null || true


# ----- Clean
.PHONY: clean

clean: test.clean ## Clean all artifacts

# ----- Help
.PHONY: help

help: ## Show help information
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST);

print-%: ; @echo "$($*)"
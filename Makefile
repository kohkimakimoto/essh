.PHONY:help dev dist packaging fmt test testv deps deps_update website
.DEFAULT_GOAL := help

# This is a magic code to output help message at default
# see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: ## Build dev binary
	@bash -c $(CURDIR)/_build/dev.sh

dist: ## Build dist binaries 
	@bash -c $(CURDIR)/_build/dist.sh

packaging: ## Create packages (now support RPM only)
	@bash -c $(CURDIR)/_build/packaging.sh

fmt: ## Run `go fmt`
	go fmt $$(go list ./... | grep -v vendor)

test: ## Run tests
	@bash -c $(CURDIR)/test/test.sh

testv: ## Run tests with verbose outputting.
	@export GOTEST_FLAGS="-cover -timeout=360s -v" && bash -c $(CURDIR)/test/test.sh

deps: ## Install dependences by using glide
	glide install

deps_update:  ## Update dependences by using glide
	glide up

website: ## Build webside resources.
	cd website && make deps && make

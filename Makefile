.PHONY: default dev dist packaging fmt test testv deps deps_update website

default: dev

dev:
	@bash -c $(CURDIR)/_build/dev.sh

dist:
	@bash -c $(CURDIR)/_build/dist.sh

packaging:
	@bash -c $(CURDIR)/_build/packaging.sh

fmt:
	go fmt $$(go list ./... | grep -v vendor)

test:
	@bash -c $(CURDIR)/test/test.sh

testv:
	@export GOTEST_FLAGS="-cover -timeout=360s -v" && bash -c $(CURDIR)/test/test.sh

# test_integration:
# 	cd tests && ./run.sh

deps:
	gom install

deps_update:
	rm Gomfile.lock; rm -rf vendor; gom install && gom lock

website:
	cd website && make deps && make

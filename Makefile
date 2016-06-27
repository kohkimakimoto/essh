.PHONY: default dev dist packaging packaging_destroy fmt test testv deps updatedeps

default: dev

dev:
	@bash -c $(CURDIR)/_build/dev.sh

dist:
	@bash -c $(CURDIR)/_build/dist.sh

packaging:
	@bash -c $(CURDIR)/_build/packaging.sh

packaging_destroy:
	@sh -c "cd $(CURDIR)/_build/packaging/rpm && vagrant destroy -f"

fmt:
	go fmt ./...

test:
	go test ./... -cover

testv:
	go test ./... -v

deps:
	gom install

deps_update:
	rm Gomfile.lock; rm -rf _vendor; gom install && gom lock

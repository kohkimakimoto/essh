.PHONY: default dev dist packaging test testv deps updatedeps

default: dev

# dev creates your platform binary.
dev:
	@sh -c "$(CURDIR)/build/build.sh dev"

# dist creates all platform binaries.
dist:
	@sh -c "$(CURDIR)/build/build.sh dist"

# packaging creates all platform binaries and rpm packages.
packaging:
	@sh -c "$(CURDIR)/build/build.sh packaging"

test:
	gom test ./... -cover

testv:
	gom test ./... -v

deps:
	gom install

updatedeps:
	rm Gomfile.lock; rm -rf _vendor; gom install && gom lock

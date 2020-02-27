BINDIR      := $(CURDIR)/bin
DIST_DIRS   := find * -maxdepth 0 -type d -exec
TARGETS     := darwin/amd64 linux/amd64

GOPATH        = $(shell go env GOPATH)
GOX           = $(GOPATH)/bin/gox

TAGS       :=
LDFLAGS    := -w -s
GOFLAGS    :=
SRC        := $(shell find . -type f -name '*.go' -print)
OBJS       := $(foreach name,$(wildcard cmd/*), $(subst cmd/,$(BINDIR)/k8s-,$(name)))

SHELL      = /bin/bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

.PHONY: all
all: build

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(OBJS)

$(OBJS): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o "$@" "./cmd/$(subst $(BINDIR)/k8s-,,$@)"

# ------------------------------------------------------------------------------
#  dependencies

$(GOX):
	(cd /; GO111MODULE=on go get -u github.com/mitchellh/gox)

# ------------------------------------------------------------------------------
#  release

define build-cross
	GO111MODULE=on CGO_ENABLED=0 $(GOX) -parallel=3 -output="_dist/{{.OS}}-{{.Arch}}/bin/k8s-$1" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/$1;
endef

.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross: $(GOX)
	$(foreach name,$(wildcard cmd/*), $(call build-cross,$(subst cmd/,,$(name))))

.PHONY: dist
dist: build-cross
	( \
		cd _dist && \
		$(DIST_DIRS) cp ../LICENSE {} \; && \
		$(DIST_DIRS) cp ../README.md {} \; && \
		$(DIST_DIRS) tar -zcf k8s-busybox-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r k8s-busybox-${VERSION}-{}.zip {} \; \
	)

.PHONY: checksum
checksum:
	for f in _dist/*.{gz,zip} ; do \
		shasum -a 256 "$${f}"  | awk '{print $$1}' > "$${f}.sha256" ; \
	done

# ------------------------------------------------------------------------------

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist

.PHONY: info
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"

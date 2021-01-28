BIN = goxlr-routing
BUILDOS ?= windows
BUILDARCH ?= 386
BUILDENVS ?= CGO_ENABLED=0 GOOS=$(BUILDOS) GOARCH=$(BUILDARCH)
BUILDFLAGS ?= -a -installsuffix cgo --ldflags '-X main.Version=$(VERSION) -extldflags "-lm -lstdc++ -static"'

build:
	@echo "==> Building binary ($(BUILDOS)/$(BUILDARCH)/$(BIN))"
	@$(BUILDENVS) go build -v $(BUILDFLAGS) -o bin/$(BIN).exe

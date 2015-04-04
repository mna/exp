SHELL = /bin/sh

ROOT = .
ROOTDIR = $(ROOT)/*.go
BINDIR = ./bin

BUILDER_DIR = $(ROOT)/builder
BUILDER_SRC = $(BUILDER_DIR)/*.go
AST_DIR = $(ROOT)/ast
AST_SRC = $(AST_DIR)/*.go

# bootstrap tools variables
BOOTSTRAP_DIR = $(ROOT)/bootstrap
BOOTSTRAP_SRC = $(BOOTSTRAP_DIR)/*.go
BOOTSTRAPBUILD_DIR = $(BOOTSTRAP_DIR)/cmd/bootstrap-build
BOOTSTRAPBUILD_SRC = $(BOOTSTRAPBUILD_DIR)/*.go
BOOTSTRAPPIGEON_DIR = $(BOOTSTRAP_DIR)/cmd/bootstrap-pigeon
BOOTSTRAPPIGEON_SRC = $(BOOTSTRAPPIGEON_DIR)/*.go

# grammar variables
GRAMMAR_DIR = $(ROOT)/grammar
BOOTSTRAP_GRAMMAR = $(GRAMMAR_DIR)/bootstrap.peg
PIGEON_GRAMMAR = $(GRAMMAR_DIR)/pigeon.peg

all: $(BINDIR)/bootstrap-build $(BOOTSTRAPPIGEON_DIR)/bootstrap_pigeon.go \
	$(BINDIR)/bootstrap-pigeon $(ROOT)/pigeon.go

$(BINDIR)/bootstrap-build: $(BOOTSTRAPBUILD_SRC) $(BOOTSTRAP_SRC) $(BUILDER_SRC) \
	$(AST_SRC)
	go build -o $@ $(BOOTSTRAPBUILD_DIR)

$(BOOTSTRAPPIGEON_DIR)/bootstrap_pigeon.go: $(BINDIR)/bootstrap-build \
	$(BOOTSTRAP_GRAMMAR)
	$(BINDIR)/bootstrap-build $(BOOTSTRAP_GRAMMAR) | goimports > $@

$(BINDIR)/bootstrap-pigeon: $(BOOTSTRAPPIGEON_SRC) \
	$(BOOTSTRAPPIGEON_DIR)/bootstrap_pigeon.go
	go build -o $@ $(BOOTSTRAPPIGEON_DIR)

$(ROOT)/pigeon.go: $(BINDIR)/bootstrap-pigeon $(PIGEON_GRAMMAR)
	$(BINDIR)/bootstrap-pigeon $(PIGEON_GRAMMAR) | goimports > $@

$(BINDIR)/pigeon: $(ROOT_SRC) $(ROOT)/pigeon.go
	go build -o $@ $(ROOT)

$(BOOTSTRAP_GRAMMAR):
$(PIGEON_GRAMMAR):

clean:
	rm $(BOOTSTRAPBUILD_DIR)/bootstrap-build $(BOOTSTRAPPIGEON_DIR)/bootstrap-pigeon \
		$(BOOTSTRAPPIGEON_DIR)/bootstrap_pigeon.go $(ROOT)/pigeon.go

.PHONY: all clean


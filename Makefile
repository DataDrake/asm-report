.RECIPEPREFIX != ps

include Makefile.waterlog

GOCC     = go

GOPATH   = $(shell pwd)/build
GOBIN    = build/bin
GOSRC    = build/src
PROJROOT = $(GOSRC)/github.com/DataDrake
PROJNAME = asm-report

DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

all: build

build: setup
    @$(call stage,BUILD)
    @$(GOCC) install -v github.com/DataDrake/$(PROJNAME)
    @$(call pass,BUILD)

setup:
    @$(call stage,SETUP)
    @$(call task,Setting up GOPATH...)
    @mkdir -p $(GOPATH)
    @$(call task,Setting up src/...)
    @mkdir -p $(GOSRC)
    @$(call task,Setting up project root...)
    @mkdir -p $(PROJROOT)
    @$(call task,Setting up symlinks...)
    @if [ ! -d $(PROJROOT)/$(PROJNAME) ]; then ln -s $(shell pwd) $(PROJROOT)/$(PROJNAME); fi
    @$(call task,Getting dependencies...)
    @go get github.com/DataDrake/waterlog
    @$(call pass,SETUP)

validate: golint-setup
    @$(call stage,FORMAT)
    @$(GOCC) fmt -x ./...
    @$(call pass,FORMAT)
    @$(call stage,VET)
    @$(GOCC) vet -x ./...
    @$(call pass,VET)
    @$(call stage,LINT)
    @$(GOBIN)/golint -set_exit_status ./...
    @$(call pass,LINT)

golint-setup:
    @if [ ! -e $(GOBIN)/golint ]; then \
        printf "Installing golint..."; \
        $(GOCC) get -u github.com/golang/lint/golint; \
        printf "DONE\n\n"; \
        rm -rf $(GOPATH)/src/golang.org $(GOPATH)/src/github.com/golang $(GOPATH)/pkg; \
    fi

install:
    @$(call stage,INSTALL)
    install -D -m 00755 $(GOBIN)/$(PROJNAME) $(DESTDIR)$(BINDIR)/$(PROJNAME)
    @$(call pass,INSTALL)

uninstall:
    @$(call stage,UNINSTALL)
    rm -f $(DESTDIR)$(BINDIR)/$(PROJNAME)
    @$(call pass,UNINSTALL)

clean:
    @$(call stage,CLEAN)
    @$(call task,Removing symlinks...)
    @unlink $(PROJROOT)/$(PROJNAME)
    @$(call task,Removing build directory...)
    @rm -rf build
    @$(call pass,CLEAN)

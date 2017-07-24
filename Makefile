# Makefile for Board project
#
# Targets:
#   all: Builds the code
#   build: Builds the code
#   fmt: Formats the source files
#   clean_binary: cleans the code
#   install: Installs the code to the GOPATH
#   test: Runs the tests
#   vet: Vet examines source code and reports suspicious constructs
#   golint: Linter for source code
#
#

# Common
BASEIMAGE=ubuntu:14.04
GOBUILDIMAGE=golang:1.8.1

# Base shell parameters
SHELL := /bin/bash
BUILDPATH=$(CURDIR)
MAKEPATH=$(BUILDPATH)/make
MAKEDEVPATH=$(MAKEPATH)/dev
SRCPATH= src
TOOLSPATH=$(BUILDPATH)/tools

# docker parameters
DOCKERCMD=$(shell which docker)
DOCKERBUILD=$(DOCKERCMD) build
DOCKERRMIMAGE=$(DOCKERCMD) rmi
DOCKERPULL=$(DOCKERCMD) pull
DOCKERIMASES=$(DOCKERCMD) images
DOCKERSAVE=$(DOCKERCMD) save
DOCKERCOMPOSECMD=$(shell which docker-compose)
DOCKERTAG=$(DOCKERCMD) tag

DOCKERCOMPOSEFILEPATH=$(MAKEDEVPATH)
DOCKERCOMPOSEFILENAME=docker-compose.yml

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
#GODEP=$(GOTEST) -i
GOFMT=gofmt -w
GOVET=$(GOCMD) vet
GOLINT=golint

# prepare parameters
PREPAREPATH=$(TOOLSPATH)
PREPARECMD=prepare
PREPARECMD_PARAMETERS=--conf $(CONFIGPATH)/$(CONFIGFILE)

# Package lists
TOPLEVEL_PKG := .
INT_LIST := $(SRCPATH)/apiserver $(SRCPATH)/tokenserver $(SRCPATH)/collector/cmd

# List building
ALL_LIST = $(INT_LIST) #$(IMPL_LIST) $(CMD_LIST)

BUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_build)
CLEAN_LIST = $(foreach int, $(ALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(ALL_LIST), $(int)_install)
TEST_LIST = $(foreach int, $(ALL_LIST), $(int)_test)
FMT_LIST = $(foreach int, $(ALL_LIST), $(int)_fmt)
VET_LIST = $(foreach int, $(ALL_LIST), $(int)_vet)
GOLINT_LIST = $(foreach int, $(ALL_LIST), $(int)_golint)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST) $(VET_LIST) $(GOLINT_LIST)

all: build
build: $(BUILD_LIST)
clean_binary: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
fmt: $(FMT_LIST)
vet: $(VET_LIST)
golint: $(GOLINT_LIST)

$(BUILD_LIST): %_build: %_fmt %_vet %_golint
	cd $(TOPLEVEL_PKG)/$*/; $(GOBUILD) .
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(TOPLEVEL_PKG)/$* 
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(TEST_LIST): %_test:
	$(GOTEST) $(TOPLEVEL_PKG)/$*
$(FMT_LIST): %_fmt:
	$(GOFMT) ./$*
$(VET_LIST): %_vet:
	$(GOVET) ./$*/...
$(GOLINT_LIST): %_golint:
	$(GOLINT) $*/...

prepare:
	@echo "preparing..."
	@$(MAKEPATH)/$(PREPARECMD) $(PREPARECMD_PARA)

start:
	@echo "loading Board images..."
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAME) up -d
	@echo "Start complete. You can visit Board now."

down:
	@echo "stoping Board instance..."
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAME) down -v
	@echo "Done."
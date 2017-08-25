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
DEVIMAGEPATH= make/dev

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
DOCKERCOMPOSEUIFILENAME=docker-compose.uibuilder.yml

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

# swagger parameters
SWAGGERTOOLPATH=$(TOOLSPATH)/swagger
SWAGGERFILEPATH=$(BUILDPATH)/docs

# Package lists
TOPLEVEL_PKG := .
INT_LIST := apiserver tokenserver collector/cmd
IMG_LIST := apiserver tokenserver log collector 

# List building
COMPILEALL_LIST = $(foreach int, $(INT_LIST), $(SRCPATH)/$(int))

COMPILE_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_compile)
CLEAN_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_install)
TEST_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_test)
FMT_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_fmt)
VET_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_vet)
GOLINT_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_golint)

BUILDALL_LIST = $(foreach int, $(IMG_LIST), container/$(int))
BUILD_LIST = $(foreach int, $(BUILDALL_LIST), $(int)_build)
RMIMG_LIST = $(foreach int, $(BUILDALL_LIST), $(int)_rmi)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(COMPILE_LIST) $(VET_LIST) $(GOLINT_LIST) $(BUILD_LIST)

all: compile 
compile: $(COMPILE_LIST) compile_ui
cleanbinary: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
fmt: $(FMT_LIST)
vet: $(VET_LIST)
golint: $(GOLINT_LIST)

compile_ui:
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEUIFILENAME) up

$(COMPILE_LIST): %_compile: %_fmt %_vet %_golint
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

build: $(BUILD_LIST) container/db_build
cleanimage: $(RMIMG_LIST) container/db_rmi

$(BUILD_LIST): %_build:
	$(DOCKERBUILD) -f $(MAKEDEVPATH)/$*/Dockerfile . -t dev_$(subst container/,,$*):latest
$(RMIMG_LIST): %_rmi:
	$(DOCKERRMIMAGE) dev_$(subst container/,,$*):latest

container/db_build:
	$(DOCKERBUILD) -f $(MAKEDEVPATH)/container/db/Dockerfile . -t dev_mysql:latest
container/db_rmi:
	$(DOCKERRMIMAGE) dev_mysql:latest

prepare:
	@echo "preparing..."
	@$(MAKEPATH)/$(PREPARECMD) $(PREPARECMD_PARA)
	@echo "Done."

start:
	@echo "loading Board images..."
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAME) up -d
	@echo "Start complete. You can visit Board now."

down:
	@echo "stoping Board instance..."
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAME) down -v
	@echo "Done."

prepare_swagger:
	@echo "preparing swagger environment..."
	@cd $(SWAGGERTOOLPATH); ./prepare-swagger.sh
	@echo "Done."

.PHONY: cleanall
cleanall: cleanbinary cleanimage

clean:
	@echo "  make cleanall:         remove binaries and Board images"
	@echo "  make cleanbinary:      remove apiserver tokenserver and collector/cmd binary"
	@echo "  make cleanimage:       remove Board images"

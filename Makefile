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
# Develop flag
#
DEVFLAG=release

# ARCH default is x86_64, also support mips, arm64v8
ARCH=

ifeq ($(DEVFLAG), release) 
	BASEIMAGE=alpine:3.7
	GOBUILDIMAGE=inspursoft/go-builder:1.14.1-alpine
	WORKPATH=release
	IMAGEPREFIX=board
else
	BASEIMAGE=ubuntu:14.04
	GOBUILDIMAGE=golang:1.14.0
	WORKPATH=dev
	IMAGEPREFIX=dev
endif 

ifeq ($(ARCH), mips)
	GOBUILDIMAGE=inspursoft/golang-mips:1.12.9
endif

ifeq ($(ARCH), arm64v8)
        GOBUILDIMAGE=inspursoft/go-builder-arm64v8:1.13.5-alpine
endif

# Base shell parameters
SHELL := /bin/bash
BUILDPATH=$(CURDIR)
MAKEPATH=$(BUILDPATH)/make
MAKEWORKPATH=$(MAKEPATH)/$(WORKPATH)
SRCPATH= src
TOOLSPATH=$(BUILDPATH)/tools
IMAGEPATH=$(BUILDPATH)/make/$(MAKEWORKPATH)
PACKAGEPATH=$(BUILDPATH)/Deploy


# docker parameters
DOCKERCMD=$(shell which docker)
DOCKERBUILD=$(DOCKERCMD) build
DOCKERRMIMAGE=$(DOCKERCMD) rmi
DOCKERPULL=$(DOCKERCMD) pull
DOCKERIMASES=$(DOCKERCMD) images
DOCKERSAVE=$(DOCKERCMD) save
DOCKERCOMPOSECMD=$(shell which docker-compose)
DOCKERTAG=$(DOCKERCMD) tag
DOCKERNETWORK=$(DOCKERCMD) network

DOCKERCOMPOSEFILEPATH=$(MAKEWORKPATH)
DOCKERCOMPOSEFILENAME=docker-compose${if ${ARCH},.${ARCH}}.yml
DOCKERCOMPOSEFILENAMEADM=docker-compose-adminserver${if ${ARCH},.${ARCH}}.yml
DOCKERCOMPOSEUIFILENAME=docker-compose.uibuilder${if ${ARCH},.${ARCH}}.yml

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
GOIMGBASEPATH=/go/src/git/inspursoft/board

# prepare parameters
PREPAREPATH=$(TOOLSPATH)
PREPARECMD=prepare
PREPARECMD_PARAMETERS=--conf $(CONFIGPATH)/$(CONFIGFILE)

# swagger parameters
SWAGGERTOOLPATH=$(TOOLSPATH)/swagger
SWAGGERFILEPATH=$(BUILDPATH)/docs

#package 
TARCMD=$(shell which tar)
ZIPCMD=$(shell which gzip)
PKGTEMPPATH=Deploy
PKGNAME=board
GITTAGVERSION=$(shell git describe --tags || echo UNKNOWN)
VERSIONFILE=VERSION
ifeq ($(DEVFLAG), release)
	VERSIONTAG=$(GITTAGVERSION)
else
	VERSIONTAG=dev
endif

# Package lists
# TOPLEVEL_PKG := .
INT_LIST := adminserver apiserver tokenserver
ifndef ARCH
	IMG_LIST := adminserver apiserver tokenserver log jenkins db proxy proxy_adminserver gogits grafana elasticsearch kibana chartmuseum prometheus
else
	IMG_LIST := apiserver tokenserver log jenkins db proxy gogits prometheus
endif


# List building
COMPILEALL_LIST = $(foreach int, $(INT_LIST), $(SRCPATH)/$(int))

COMPILE_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_compile)
CLEAN_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_install)
TEST_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_test)
FMT_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_fmt)
VET_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_vet)
GOLINT_LIST = $(foreach int, $(COMPILEALL_LIST), $(int)_golint)
PKG_LIST = $(foreach int, $(IMG_LIST), $(IMAGEPREFIX)_$(int):$(VERSIONTAG))

BUILDALL_LIST = $(foreach int, $(IMG_LIST), container/$(int))
BUILD_LIST = $(foreach int, $(BUILDALL_LIST), $(int)_build)
RMIMG_LIST = $(foreach int, $(BUILDALL_LIST), $(int)_rmi)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(COMPILE_LIST) $(VET_LIST) $(GOLINT_LIST) $(BUILD_LIST)

all: compile 
compile: $(COMPILE_LIST)
cleanbinary: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
fmt: $(FMT_LIST)
vet: $(VET_LIST)
golint: $(GOLINT_LIST)

version:
	@echo $(VERSIONTAG)
	@echo $(VERSIONTAG) > $(VERSIONFILE)

compile_ui:
	$(DOCKERCOMPOSECMD) -f $(MAKEWORKPATH)/$(DOCKERCOMPOSEUIFILENAME) up
	$(DOCKERCOMPOSECMD) -f $(MAKEWORKPATH)/$(DOCKERCOMPOSEUIFILENAME) down

$(COMPILE_LIST): %_compile: # %_fmt  %_vet %_golint
	$(DOCKERCMD) run --rm -v $(BUILDPATH):$(GOIMGBASEPATH) \
					-w $(GOIMGBASEPATH)/$* $(GOBUILDIMAGE) $(GOBUILD) \
					-v -o $(GOIMGBASEPATH)/make/$(WORKPATH)/container/$(subst /cmd,,$(subst src/,,$*))/$(subst /cmd,,$(subst src/,,$*)) 

$(CLEAN_LIST): %_clean:
#	$(GOCLEAN) $(TOPLEVEL_PKG)/$* 
	rm $(MAKEWORKPATH)/container/$(subst /cmd,,$(subst src/,,$*))/$(subst /cmd,,$(subst src/,,$*))	    
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(TEST_LIST): %_test:
#	$(GOTEST) $(TOPLEVEL_PKG)/$*
	$(DOCKERCMD) run --rm -v $(BUILDPATH):$(GOIMGBASEPATH) \
                                        -w $(GOIMGBASEPATH)/$* $(GOBUILDIMAGE) $(GOTEST) \
                                        -v -o $(GOIMGBASEPATH)/make/$(WORKPATH)/container/$(subst /cmd,,$(subst src/,,$*))/$(subst /cmd,,$(subst src/,,$*))
$(FMT_LIST): %_fmt:
	$(GOFMT) ./$*
$(VET_LIST): %_vet:
	$(GOVET) ./$*/...
$(GOLINT_LIST): %_golint:
	$(GOLINT) $*/...

build: version $(BUILD_LIST) #container/db_build
cleanimage: $(RMIMG_LIST) #container/db_rmi

$(BUILD_LIST): %_build: 
	$(DOCKERBUILD) -f $(MAKEWORKPATH)/$*/Dockerfile${if ${ARCH},.${ARCH}} . -t $(IMAGEPREFIX)_$(subst container/,,$*):$(VERSIONTAG)
	
$(RMIMG_LIST): %_rmi:
	$(DOCKERRMIMAGE) -f $(IMAGEPREFIX)_$(subst container/,,$*):$(VERSIONTAG)

#container/db_build:
#	$(DOCKERBUILD) -f $(MAKEWORKPATH)/container/db/Dockerfile . -t $(IMAGEPREFIX)_mysql:latest
#container/db_rmi:
#	$(DOCKERRMIMAGE) $(IMAGEPREFIX)_mysql:latest

prepare: version
	@echo "preparing..."
	@$(MAKEPATH)/$(PREPARECMD) $(PREPARECMD_PARA)
	@echo "Done."

start:
	@echo "loading Board images..."
	$(DOCKERNETWORK) create board &> /dev/null || true
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAME) up -d
	@echo "Start complete. You can visit Board now."

start_legacy:
	@echo "loading Board images..."
	$(DOCKERNETWORK) create board &> /dev/null || true
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/archive/$(DOCKERCOMPOSEFILENAME) up -d
	@echo "Start complete. You can visit Board now."

start_admin:
	@echo "loading Adminserver images..."
	@if [ ! -d $(MAKEPATH)/adminserver ] ; then mkdir -p $(MAKEPATH)/adminserver ; fi
	@rm -f $(MAKEPATH)/adminserver/env
	@cp $(MAKEPATH)/templates/adminserver/env-dev $(MAKEPATH)/adminserver/env
	@sed -i "s|__CURDIR__|$(MAKEPATH)|g"  $(MAKEPATH)/adminserver/env
	@if [ ! -f $(MAKEPATH)/adminserver/board.cfg ] ; then cp $(MAKEPATH)/board.cfg $(MAKEPATH)/adminserver/board.cfg ; fi
	$(DOCKERNETWORK) create board &> /dev/null || true
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAMEADM) up -d
	@echo "Start complete. You can visit Adminserver now."

down:
	@echo "stoping Board instance..."
	$(DOCKERNETWORK) rm board &> /dev/null || true
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAME) down -v
	@echo "Done."

down_legacy:
	@echo "stoping Board instance..."
	$(DOCKERNETWORK) rm board &> /dev/null || true
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/archive/$(DOCKERCOMPOSEFILENAME) down -v
	@echo "Done."

down_admin:
	@echo "stoping Adminserver instance..."
	$(DOCKERNETWORK) rm board &> /dev/null || true
	$(DOCKERCOMPOSECMD) -f $(DOCKERCOMPOSEFILEPATH)/$(DOCKERCOMPOSEFILENAMEADM) down -v
	@echo "Done."

prepare_swagger:
	@echo "preparing swagger environment..."
	@cd $(SWAGGERTOOLPATH); ./prepare-swagger.sh
	@echo "Done."

prepare_composefile:
	@cp $(MAKEWORKPATH)/docker-compose${if ${ARCH},.${ARCH}}.tpl $(MAKEWORKPATH)/docker-compose${if ${ARCH},.${ARCH}}.yml
	@cp $(MAKEWORKPATH)/archive/docker-compose${if ${ARCH},.${ARCH}}.tpl $(MAKEWORKPATH)/archive/docker-compose${if ${ARCH},.${ARCH}}.yml
	@cp $(MAKEWORKPATH)/docker-compose-adminserver${if ${ARCH},.${ARCH}}.tpl $(MAKEWORKPATH)/docker-compose-adminserver${if ${ARCH},.${ARCH}}.yml
	@sed -i "s/__version__/$(VERSIONTAG)/g" $(MAKEWORKPATH)/docker-compose${if ${ARCH},.${ARCH}}.yml
	@sed -i "s/__version__/$(VERSIONTAG)/g" $(MAKEWORKPATH)/archive/docker-compose${if ${ARCH},.${ARCH}}.yml
	@sed -i "s/__version__/$(VERSIONTAG)/g" $(MAKEWORKPATH)/docker-compose-adminserver${if ${ARCH},.${ARCH}}.yml

package: prepare_composefile
	@echo "packing offline package ..."
	@if [ ! -d $(PKGTEMPPATH) ] ; then mkdir $(PKGTEMPPATH) ; fi
	@if [ ! -d $(PKGTEMPPATH)/adminserver ] ; then mkdir -p $(PKGTEMPPATH)/adminserver ; fi
	@if [ ! -d $(PKGTEMPPATH)/archive ] ; then mkdir -p $(PKGTEMPPATH)/archive ; fi
	@cp $(TOOLSPATH)/install.sh $(PKGTEMPPATH)/install.sh
	@cp $(TOOLSPATH)/install-adminserver.sh $(PKGTEMPPATH)/adminserver/install-adminserver.sh
	@cp $(TOOLSPATH)/uninstall.sh $(PKGTEMPPATH)/uninstall.sh
	@cp $(TOOLSPATH)/uninstall-adminserver.sh $(PKGTEMPPATH)/adminserver/uninstall-adminserver.sh
	@cp $(MAKEPATH)/board.cfg $(PKGTEMPPATH)/.
	@cp $(MAKEPATH)/board.cfg $(PKGTEMPPATH)/adminserver/.
	@cp $(MAKEPATH)/prepare $(PKGTEMPPATH)/.
	@cp -rf $(MAKEPATH)/templates $(PKGTEMPPATH)/.
	@cp $(MAKEWORKPATH)/docker-compose${if ${ARCH},.${ARCH}}.yml $(PKGTEMPPATH)/docker-compose.yml
	@cp $(MAKEWORKPATH)/archive/docker-compose${if ${ARCH},.${ARCH}}.yml $(PKGTEMPPATH)/archive/docker-compose.yml
	@cp $(MAKEWORKPATH)/docker-compose-adminserver${if ${ARCH},.${ARCH}}.yml $(PKGTEMPPATH)/adminserver/docker-compose-adminserver.yml
	@cp $(MAKEPATH)/templates/adminserver/env-release $(PKGTEMPPATH)/adminserver/env
#	@cp LICENSE $(PKGTEMPPATH)/LICENSE
#	@cp NOTICE $(PKGTEMPPATH)/NOTICE
	@sed -i "s/..\/config/.\/config/" $(PKGTEMPPATH)/docker-compose.yml
	@sed -i "s/..\/config/.\/config/" $(PKGTEMPPATH)/archive/docker-compose.yml

	@echo "package images ..."
	@$(DOCKERSAVE) -o $(PKGTEMPPATH)/$(IMAGEPREFIX)_deployment.$(VERSIONTAG).tgz $(PKG_LIST) k8s_install:1.18 gitlab-helper:1.0
	@$(TARCMD) -zcvf $(PKGNAME)-offline-installer-$(VERSIONTAG)${if ${ARCH},.${ARCH}}.tgz $(PKGTEMPPATH)

	@rm -rf $(PACKAGEPATH)

packageonestep: compile compile_ui build package
#packageonestep: compile build package

.PHONY: cleanall
cleanall: cleanbinary cleanimage

clean:
	@echo "  make cleanall:         remove binaries and Board images"
	@echo "  make cleanbinary:      remove apiserver tokenserver binary"
	@echo "  make cleanimage:       remove Board images"

# Makefile for Board project
#
# Targets:
#
# all:			start
# compile: 		compile apiserver, tokenserver code
#
# compile_apiserver, compile_tokenserver: compile specific binary
#
# clean:        remove binary and images
# cleanbinary:	remove apiserver, tokenserver

# Base shell parameters
SHELL := /bin/bash
BUILDPATH=$(CURDIR)
MAKEPATH=$(BUILDPATH)/make
MAKEDEVPATH=$(MAKEPATH)/dev
SRCPATH=./src
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

# go parameters
GOBASEPATH=/go/src/git
GOCMD=$(shell which go)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w
GOVET=$(GOCMD) vet
GOLINT=golint

# Common
BASEIMAGE=ubuntu:14.04
GOBUILDIMAGE=golang:1.8.1
GOBUILDPATH=$(GOBASEPATH)/inspursoft/board
GOIMAGEBUILD=$(GOCMD) build
GOBUILDMAKEPATH=$(GOBUILDPATH)/make

# Board Component
## UI

## API Server
GOBUILDPATH_APISERVER=$(GOBUILDPATH)/src/apiserver
APISERVERSOURCECODE=$(SRCPATH)/apiserver
APISERVERBINARYPATH=$(MAKEDEVPATH)/apiserver
APISERVERBINARYNAME=apiserver
DOCKERIMAGENAME_APISERVER=apiserver

## Token Server
GOBUILDPATH_TOKENSERVER=$(GOBUILDPATH)/src/tokenserver
TOKENSERVERSOURCECODE=$(SRCPATH)/tokenserver
TOKENSERVERBINARYPATH=$(MAKEDEVPATH)/tokenserver
TOKENSERVERBINARYNAME=tokenserver
DOCKERIMAGENAME_TOKENSERVER=tokenserver

## Database
DOCKERIMAGENAME_DATABASE=mysql

## Log
DOCKERIMAGENAME_LOG=log

## Nginx
DOCKERIMAGENAME_NGINX=nginx

# configfile
CONFIGPATH=$(MAKEPATH)
CONFIGFILE=board.cfg

# prepare parameters
PREPAREPATH=$(TOOLSPATH)
PREPARECMD=prepare
PREPARECMD_PARAMETERS=--conf $(CONFIGPATH)/$(CONFIGFILE)

# package
TARCMD=$(shell which tar)
ZIPCMD=$(shell which gzip)
DOCKERIMGFILE=Board

# commands
compile_apiserver:
	@echo "compiling binary for apiserver..."
	@$(GOBUILD) -o $(APISERVERBINARYPATH)/$(APISERVERBINARYNAME) $(APISERVERSOURCECODE)
	@echo "Done."

compile_tokenserver:
	@echo "compiling binary for tokenserver..."
	@$(GOBUILD) -o $(TOKENSERVERBINARYPATH)/$(TOKENSERVERBINARYNAME) $(TOKENSERVERSOURCECODE)
	@echo "Done."

compile_normal: compile_apiserver compile_tokenserver

compile_golangimage: 
	@echo "compiling binary for apiserver (golang image)..."
	@echo $(GOBASEPATH)
	@echo $(GOBUILDPATH)
	@$(DOCKERCMD) run --rm -v $(BUILDPATH):$(GOBUILDPATH) -w $(GOBUILDPATH_APISERVER) $(GOBUILDIMAGE) $(GOIMAGEBUILD) -v -o $(GOBUILDMAKEPATH)/dev/$(APISERVERBINARYNAME)/$(APISERVERBINARYNAME) # <-need improve
	@echo "Done."

	@echo "compiling binary for tokenserver (golang image)..."
	@echo $(GOBASEPATH)
	@echo $(GOBUILDPATH)
	@$(DOCKERCMD) run --rm -v $(BUILDPATH):$(GOBUILDPATH) -w $(GOBUILDPATH_TOKENSERVER) $(GOBUILDIMAGE) $(GOIMAGEBUILD) -v -o $(GOBUILDMAKEPATH)/dev/$(TOKENSERVERBINARYNAME)/$(TOKENSERVERBINARYNAME) # <-need improve
	@echo "Done."

compile: compile_golangimage
	
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

cleanbinary:
	@echo "cleaning binary..."
	@if [ -f $(APISERVERBINARYPATH)/$(APISERVERBINARYNAME) ] ; then rm $(APISERVERBINARYPATH)/$(APISERVERBINARYNAME) ; fi
	@if [ -f $(TOKENSERVERBINARYPATH)/$(TOKENSERVERBINARYNAME) ] ; then rm $(TOKENSERVERBINARYPATH)/$(TOKENSERVERBINARYNAME) ; fi

cleanimage:
	@echo "cleaning images for Board..."

.PHONY: cleanall
cleanall: cleanbinary cleanimage 

clean:
	@echo "  make cleanall:		        remove binaries and Board images"
	@echo "  make cleanbinary:		remove apiserver and tokenserver binary"
	@echo "  make cleanimage:		remove Board images"

all: start

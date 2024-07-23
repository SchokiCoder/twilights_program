# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

APP_NAME:=twilights_program

INSTALLDIR_ASSETS:=/usr/share/
INSTALLDIR_ASSETS:=$(INSTALLDIR_ASSETS)$(APP_NAME)
INSTALLDIR_BIN:=/usr/local/bin

# uncomment for user install
#INSTALLDIR_ASSETS:=$(HOME)/.local/share/
#INSTALLDIR_ASSETS:=$(INSTALLDIR_ASSETS)$(APP_NAME)
#INSTALLDIR_BIN:=$(HOME)/.local/bin

.PHONY: all build clean vet install uninstall

all: vet build

build:
	go build

clean:
	rm -f $(APP_NAME)

vet:
	go vet

install: build
	cp $(APP_NAME) $(INSTALLDIR_BIN)/
	mkdir $(INSTALLDIR_ASSETS)
	cp -r assets/* $(INSTALLDIR_ASSETS)/

uninstall:
	rm $(INSTALLDIR_BIN)/$(APP_NAME)
	rm -rf $(INSTALLDIR_ASSETS)/*
	rmdir $(INSTALLDIR_ASSETS)

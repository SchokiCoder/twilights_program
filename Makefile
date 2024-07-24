# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

APP_NAME         :=twilights_program
LICENSE          :=GPL-2.0-or-later
LICENSE_URL      :=https://www.gnu.org/licenses/gpl-2.0.html
REPOSITORY       :=https://github.com/SchokiCoder/twilights_program
VERSION          :=v0.0
PATH_ASSETS_SYS  :=/usr/share/$(APP_NAME)
PATH_ASSETS_USER :=$(HOME)/.local/share/$(APP_NAME)

INSTALLDIR_ASSETS :=$(PATH_ASSETS_SYS)
INSTALLDIR_BIN    :=/usr/local/bin

# uncomment for user install
#INSTALLDIR_ASSETS :=$(PATH_ASSETS_USER)
#INSTALLDIR_BIN    :=$(HOME)/.local/bin

.PHONY: all build clean vet install uninstall

all: vet build

build:
	go build \
		-ldflags "-X 'main.AppName=$(APP_NAME)' -X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppRepository=$(REPOSITORY)' -X 'main.AppVersion=$(VERSION)' -X 'main.PathAssetsSys=$(PATH_ASSETS_SYS)' -X 'main.PathAssetsUser=$(PATH_ASSETS_USER)'"

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

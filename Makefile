# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

APP_NAME         :=twilights_program
LICENSE          :=GPL-2.0-or-later
LICENSE_URL      :=https://www.gnu.org/licenses/gpl-2.0.html
REPOSITORY       :=https://github.com/SchokiCoder/twilights_program
VERSION          :=v1.1
GO_COMPILE_VARS  :=-ldflags "-X 'main.AppName=$(APP_NAME)' -X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppRepository=$(REPOSITORY)' -X 'main.AppVersion=$(VERSION)'"
SRC              :=consts.go main.go pony_model.go sprite.go

INSTALLDIR_PARENT :=$(HOME)/.local/bin
INSTALLDIR        :=$(INSTALLDIR_PARENT)/$(APP_NAME)_data

.PHONY: all build clean vet install uninstall

all: vet build

build: $(APP_NAME)

clean:
	rm -f $(APP_NAME)
	rm -f $(APP_NAME).exe
	rm -f package_linux_amd64.tar.gz
	rm -f package_windows_amd64.zip

vet:
	go vet

install: build
	mkdir -p $(INSTALLDIR)
	cp $(APP_NAME) $(INSTALLDIR_PARENT)/
	cp -r images $(INSTALLDIR)/
	cp -r fonts $(INSTALLDIR)/
	cp -r sounds $(INSTALLDIR)/

uninstall:
	rm -rf $(INSTALLDIR)/
	rm $(INSTALLDIR_PARENT)/$(APP_NAME)

packages: package_linux_amd64.tar.gz package_windows_amd64.zip

package_linux_amd64.tar.gz: $(APP_NAME)
	tar -czf $@ $< fonts/ images/ sounds/ LICENSE

$(APP_NAME): $(SRC)
	go build $(GO_COMPILE_VARS)

package_windows_amd64.zip: $(APP_NAME).exe
	zip $@ $< images/*/* fonts/* sounds/* LICENSE

$(APP_NAME).exe: $(SRC)
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
		go build -tags static -ldflags "-s -w" $(GO_COMPILE_VARS)

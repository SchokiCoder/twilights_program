# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

APP_ID           :=io.github.SchokiCoder.twilights_program
APP_NAME         :=twilights_program
LICENSE          :=GPL-2.0-or-later
LICENSE_URL      :=https://www.gnu.org/licenses/gpl-2.0.html
REPOSITORY       :=https://github.com/SchokiCoder/twilights_program
VERSION          :=v1.1
GO_COMPILE_VARS  :=-ldflags "-X 'main.AppName=$(APP_NAME)' -X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppRepository=$(REPOSITORY)' -X 'main.AppVersion=$(VERSION)'"
SRC              :=consts.go main.go pony_model.go sprite.go

DESTDIR :=$(HOME)/.local/bin

.PHONY: all build clean vet install uninstall

all: vet build

build: $(APP_NAME)

clean:
	rm -f $(APP_NAME)
	rm -f $(APP_NAME).exe
	rm -f package_linux_amd64.tar.gz
	rm -f package_windows_amd64.zip
	rm -f $(APP_NAME)-amd64.AppImage
	rm -fr AppDir

vet:
	go vet

install: build
	mkdir -p $(DESTDIR)
	cp -t $(DESTDIR)/ \
		$(APP_NAME) $(APP_NAME).svg
	mkdir -p $(DESTDIR)/$(APP_NAME)_data/
	cp -r -t $(DESTDIR)/$(APP_NAME)_data/ images fonts sounds

uninstall:
	rm -f $(DESTDIR)/$(APP_NAME) $(DESTDIR)/$(APP_NAME).svg
	rm -fr $(DESTDIR)/$(APP_NAME)_data

packages: package_linux_amd64.tar.gz package_windows_amd64.zip

# adding metainfo/appdata spawns complaints about desktop file
$(APP_NAME)-amd64.AppImage: $(APP_NAME)
	make -e DESTDIR=AppDir install
	mv AppDir/$(APP_NAME) AppDir/AppRun
	cp AppDir/$(APP_NAME).svg AppDir/$(APP_ID).svg
	cp $(APP_ID).desktop AppDir/
	appimagetool AppDir $@

package_linux_amd64.tar.gz: $(APP_NAME)
	tar -czf $@ $< fonts/ images/ sounds/ LICENSE $(APP_NAME).svg

$(APP_NAME): $(SRC)
	go build $(GO_COMPILE_VARS)

package_windows_amd64.zip: $(APP_NAME).exe
	zip $@ $< images/*/* fonts/* sounds/* LICENSE $(APP_NAME).svg

# Doesn't work under mint 22: SDL.h not found. Use 21 (and below (maybe idk)).
$(APP_NAME).exe: $(SRC)
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
		go build -tags static -ldflags "-s -w" $(GO_COMPILE_VARS)

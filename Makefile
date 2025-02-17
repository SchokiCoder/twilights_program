# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024 - 2025  Andy Frank Schoknecht

APP_ID           :=io.github.SchokiCoder.twilights_program
APP_NAME         :=twilights_program
LICENSE          :=GPL-2.0-or-later
LICENSE_URL      :=https://www.gnu.org/licenses/gpl-2.0.html
REPOSITORY       :=https://github.com/SchokiCoder/twilights_program
VERSION          :=v1.1
GO_COMPILE_VARS  :=-ldflags "-X 'main.AppId=$(APP_ID)' -X 'main.AppName=$(APP_NAME)' -X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppRepository=$(REPOSITORY)' -X 'main.AppVersion=$(VERSION)'"
SRC              :=consts.go main.go pony_model.go sprite.go

DESTDIR      :=/usr
DESKTOP_FILE :=$(APP_ID).desktop
ICON_FILE    :=$(APP_ID).svg
METAINFO_FILE:=$(APP_ID).metainfo.xml

.PHONY: all build clean vet install uninstall

all: vet build

build: $(APP_NAME)

clean:
	rm -f $(APP_NAME)
	rm -f $(APP_NAME).exe
	rm -f *.tar.gz
	rm -f *.zip
	rm -f *.AppImage
	rm -fr AppDir

vet:
	go vet

install: build
	mkdir -p $(DESTDIR)/bin/
	cp $(APP_NAME) $(DESTDIR)/bin/
	mkdir -p $(DESTDIR)/share/$(APP_NAME)/
	cp -r -t $(DESTDIR)/share/$(APP_NAME)/ images fonts sounds
	mkdir -p $(DESTDIR)/share/icons/hicolor/scalable/apps
	cp $(ICON_FILE) $(DESTDIR)/share/icons/hicolor/scalable/apps/
	mkdir -p $(DESTDIR)/share/applications/
	cp $(DESKTOP_FILE) $(DESTDIR)/share/applications/
	mkdir -p $(DESTDIR)/share/metainfo/
	cp $(METAINFO_FILE) $(DESTDIR)/share/metainfo/

uninstall:
	rm -f $(DESTDIR)/bin/$(APP_NAME)
	rm -fr $(DESTDIR)/share/$(APP_NAME)/
	rm -f $(DESTDIR)/share/icons/hicolor/scalable/apps/$(ICON_FILE)
	rm -f $(DESTDIR)/share/applications/$(DESKTOP_FILE)
	rm -f $(DESTDIR)/share/metainfo/$(METAINFO_FILE)

packages: package_linux_amd64.tar.gz package_windows_amd64.zip $(APP_NAME)-amd64.AppImage

$(APP_NAME)-amd64.AppImage: $(APP_NAME)
	make install DESTDIR=AppDir/usr
	linuxdeploy --appdir AppDir \
		-e $(APP_NAME) -d $(DESKTOP_FILE) -i $(ICON_FILE) \
		--output appimage
	mv *.AppImage $@

package_linux_amd64.tar.gz: $(APP_NAME)
	tar -czf $@ $< images fonts sounds $(ICON_FILE) $(DESKTOP_FILE) $(METAINFO_FILE) LICENSE

$(APP_NAME): $(SRC)
	go build $(GO_COMPILE_VARS)

package_windows_amd64.zip: $(APP_NAME).exe
	zip $@ $< images/*/* fonts/* sounds/* LICENSE $(ICON_FILE)

# ATTENTION, finicky:
# Cross compile for windows will only work when SDL2 and SDL2_{ttf,mixer,...}
# are installed into the mingw dirs (eg. "/usr/x86_64-w64-mingw32")
# The necessary files (headers, dlls) may be provided by package manager,
# BUT THEY MAY ALSO NOT WORK.
# For example, on Fedora 41 I get errors about missing "X11/Xlib.h".
# In such cases, don't waste your time, like I did,
# and just grab the official release files from Github.

$(APP_NAME).exe: $(SRC)
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
		go build $(GO_COMPILE_VARS)

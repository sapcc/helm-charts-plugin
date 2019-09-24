PLUGIN_NAME := charts

.PHONY: build
build: build_linux build_mac build_windows

build_windows: export GOARCH=amd64
build_windows:
	@GOOS=windows go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/windows/amd64/helm-charts main.go  # windows

link_windows:
	@cp bin/windows/amd64/helm-charts ./bin/helm-charts

build_linux: export GOARCH=amd64
build_linux: export CGO_ENABLED=0
build_linux:
	@GOOS=linux go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/linux/amd64/helm-charts main.go  # linux

link_linux:
	@cp bin/linux/amd64/helm-charts ./bin/helm-charts

build_mac: export GOARCH=amd64
build_mac: export CGO_ENABLED=0
build_mac:
	@GOOS=darwin go build -v --ldflags="-w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
                 		-o bin/darwin/amd64/helm-charts main.go # mac osx
	@cp bin/darwin/amd64/helm-charts ./bin/helm-charts # For use w make install

link_mac:
	@cp bin/darwin/amd64/helm-charts ./bin/helm-charts

.PHONY: clean
clean:
	@git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: tree
tree:
	@tree -I vendor

.PHONY: release
release:
	@scripts/release.sh $(VERSION)

.PHONY: install
install:
	HELM_OUTDATED_DEPENDENCIES_PLUGIN_NO_INSTALL_HOOK=1 helm plugin install $(shell pwd)

.PHONY: remove
remove:
	helm plugin remove $(PLUGIN_NAME)

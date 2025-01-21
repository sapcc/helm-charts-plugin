.PHONY: release
release:
	@goreleaser $(VERSION) --clean

.PHONY: install
install:
	HELM_OUTDATED_DEPENDENCIES_PLUGIN_NO_INSTALL_HOOK=1 helm plugin install $(shell pwd)

.PHONY: remove
remove:
	helm plugin remove $(PLUGIN_NAME)

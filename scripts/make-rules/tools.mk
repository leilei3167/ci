
.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* >/dev/null;then $(MAKE) tools.install.$*;fi

.PHONY: tools.install.%
tools.install.%:
	@echo "=====>install $*"
	@$(MAKE) install.$*

.PHONY: install.gofumpt
install.gofumpt:
	@$(GO) install mvdan.cc/gofumpt@latest

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

.PHONY: install.golines
install.golines:
	@$(GO) install github.com/segmentio/golines@latest

.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.mockgen
install.mockgen:
	@$(GO) install github.com/golang/mock/mockgen@latest

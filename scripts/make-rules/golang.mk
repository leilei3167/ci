GO := go

# 应该设置更多的构建选项,使得能够构建不同的平台,版本的二进制文件,便于后续CI
.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

# TODO:可以在此设置检验覆盖率,覆盖率未达标的报错
.PHONY: go.test
go.test:
	@echo "===========> Run unit test"
	@$(GO) test -race $(ROOT_DIR)/...

.PHONY: go.test.cover
go.test.cover:
	@$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out $(ROOT_DIR)/...

.PHONY: go.build
go.build:
	@echo "===========> Building binary"
	@echo "Build is unsupported now"
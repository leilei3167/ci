# 默认执行所有,包含格式化代码、静态代码检查、单元测试、代码构建、文件清理、除了帮助功能
.DEFAULT_GOAL := all

.PHONY: all
all: tidy format lint test build

# ===============================================================================
# Build options
ROOT_PACKAGE=github.com/leilei3167/ci

# ===============================================================================
# includes 此处引入其他的makefile,commom定义各个Makefile通用的环境变量,golang提供go相关的目标
# tools提供各种工具的检查,或者安装
include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk

# ===============================================================================
# Usage 此处定义flag的用法信息,一个多行变量
define USAGE_OPTIONS

Options:
	
endef
export USAGE_OPTIONS

# ===============================================================================
# Targets 常用的目标

# format: 格式化代码(golines来优化换行以及tag对齐,goimports格式化导入,会使用gofumpt格式化代码(是gofmt的超集))
# 需要先验证其工具是否安装,未安装则会安装
.PHONY: format
format: tools.verify.golines tools.verify.gofumpt tools.verify.goimports
	@echo "========>formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=140 --reformat-tags --shorten-comments --ignore-generated .
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) gofumpt -w
	@$(GO) mod edit -fmt

# lint:对代码进行静态检查
.PHONY: lint
lint:
	@$(MAKE) go.lint 

# tidy: 运行go mod tidy
.PHONY: tidy
tidy:
	@$(GO) mod tidy

# clean: 清除output文件夹
.PHONY: clean
clean:
	@echo "=======> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)
	@-rm -v api-server logic-server task
	@echo "=======> Done."

# build: 编译
.PHONY: build
build:
	@$(MAKE) go.build

# test: 测试
.PHONY: test
test:
	@$(MAKE) go.test

# cover:
.PHONY: cover
cover:
	@$(MAKE) go.test.cover

# ===============================================================================






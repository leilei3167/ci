# common 主要用于初始化变量,必须被主Makefile第一个include
SHELL := /bin/bash

# $(MAKEFILE_LIST)是特殊变量,其返回make要处理的文件列表(当前文件一定在最后一个)
# 此处用lastword函数将最后一个文件提取出,使用dir函数从文件名序列中取出目录部分
# 值为scripts/make-rules/
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)# 如果没有定义ROOT_DIR变量,cd到mk文件的目录,返回上两级并将pwd结果给abspath函数
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)# 判断是否定义了输出目录
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif

# 定义find和xargs,便于对代码处理时,跳过第三方的代码
FIND := find . ! -path './third_party/*' ! -path './vendor/*'
XARGS := xargs --no-run-if-empty
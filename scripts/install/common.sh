#!/usr/bin/env bash

# 确定工作文件夹,导入环境变量
ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")../..

source "${ROOT_DIR}/scripts/lib/init.sh"
source "${ROOT_DIR}/scripts/install/env.sh"

# 不输入密码执行需要 root 权限的命令
function common::sudo {
  echo ${LINUX_PASSWORD} | sudo -S $1
}

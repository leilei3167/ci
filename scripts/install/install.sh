#!/usr/bin/env bash

#必须设置 LINUX_PASSWORD 环境变量才能执行此脚本

ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")../..
source "${ROOT_DIR}/scripts/install/common.sh"

# 定义安装的函数
function install::install()
{
    #1.第一步,配置linux,使其成为友好的Go开发机
    install::init_into_go_env || return 1

    #2.第二步,安装编写的程序
    #3.运行测试脚本,确保功能正常
}

# 配置linux,使其成为Go开发机
function install::init_into_go_env()
{
    #1. Linux服务器的基本配置
    install::prepare_linux || return 1
    #2. Go环境的安装和配置
    install::go || return 1

    #3. Go 开发IDE的安装和配置(可选)
    #install::vim_ide || return 1

}

# 假设是centOS(后续可增加判断适配Ubuntu或CentOS)
function install::prepare_linux()
{
    #1. 配置yum源,更换为阿里的源


    #2. 配置.bashrc,后续可在此进行zsh的安装和配置


    #3. 使用yum安装必备的依赖包


    #4. 安装高版本git


    #5. 配置git,用户名,邮箱等

    #6. 配置ssh并开启(sshd配置等)

    #7. 安装docker(如果需要)
}

# 安装go开发环境
function install::go()
{
    #1. 安装go
    install::go_command || return 1

    #2. 安装protobuf
    install::protobuf || return 1



}

function install::go_command()
{
    #1. 下载go安装包


    #2. 解压安装

    #3. 配置环境变量到 .bashrc

    #4. 初始化Go工作区

}

function install::protobuf()
{
    #1. 先检查几个命令是否存在
}

function install::vim_ide()
{
    #1. 安装vim-go

    #2. Vim的一系列Go工具
}
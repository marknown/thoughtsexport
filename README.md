# 所思文档导出工具

teambition所思文档到期，需要导出所有的所思知识库，人工迁移工作量太大，写了一个小工具。目前支持docx，html两种格式导出。

## 使用说明

1. 支持 windows、linux、macos
2. 请务必安装 Chrome 浏览器。
3. 下载目录为"当前软件目录/企业名称/知识库名称/"

## 导出流程
1. 打开一个本地网页，输入一个所思知识库地址。
2. 打开所思登录页面，需要在此而面上完成登录。
3. 获取到登录态后，开始下载所思知识库。
4. 尽量使用所思知识库所有人账号登录，以免有权限问题。

## 使用方式

### 直接下载二进制运行

下载 bin 目录下对应版本的软件到自己的电脑即可运行。

系统 | 下载地址
---|---
windows | bin/thoughts_export_win
linux | bin/thoughts_export_linux
macos | bin/thoughts_export_macos


### 源码运行（需要安装golang环境）
```
git clone https://github.com/marknown/thoughtsexport.git

cd thoughtsexport
go run main.go

注意：用此命令运行时，文档下载地址为"二进制编译目录/企业名称/知识库名称/"，非当前目录。
```

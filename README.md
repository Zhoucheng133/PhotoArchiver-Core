# PhotoArchiver Core

## 简介

![License](https://img.shields.io/badge/License-MIT-dark_green)

这是PhotoArchiver软件的一部分，主仓库见[PhotoArchiver](https://github.com/Zhoucheng133/Photo-Archiver)

## 如果你想要自行打包成动态库

使用下面的命令来生成动态库
```bash
#  macOS
go build -buildmode=c-shared -ldflags="-s -w" -o build/core.dylib
# Windows
go build -buildmode=c-shared -ldflags="-s -w" -o build/core.dll
```
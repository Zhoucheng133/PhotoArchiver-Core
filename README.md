# PhotoArchiver Core

## 简介

![License](https://img.shields.io/badge/License-MIT-dark_green)

这是PhotoArchiver软件的一部分，主仓库见[PhotoArchiver](https://github.com/Zhoucheng133/Photo-Archiver)

## 如果你想要自行打包成动态库

使用下面的命令来生成动态库
```bash
# 对于Windows系统
go build -o build/core.dll -buildmode=c-shared .
# 对于macOS系统
go build -o build/core.dylib -buildmode=c-shared .
```
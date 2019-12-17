# Go 代码格式化

## 1.功能说明
* 格式化import部分：分3段式，依次为 `标准库`、`第三方库`、`项目自身库`
* 格式化单行注释：若为 `//注释内容`，调整为 `//{空格}注释内容`
* 格式化多行注释：去除首未空行；除了首未2行外，每行格式为 `{空格}*{空格}注释内容`
* 默认只对git项目库里有修改的进行格式化

## 2.安装
```
go get -u github.com/fsgo/go_fmt
```
当前版本：v0.1 20191217

## 3.使用

### 3.1 格式化git项目里有修改的.go文件
```
$ go_fmt -w
```

### 3.2 对当前目录所有.go文件格式化
```
$ go_fmt -w ./...
```

### 3.3 对指定.go文件格式化
```
$ go_fmt -w abc.go
```

## 4.配置到git hooks(commit前自动格式化代码)
```
mkdir -p ~/.git_config/hooks/
git config --global core.hooksPath ~/.git_config/hooks/

wget https://raw.githubusercontent.com/fsgo/go_fmt/master/pre-commit -o ~/.git_config/hooks/pre-commit
chmod 777 ~/.git_config/hooks/pre-commit
```
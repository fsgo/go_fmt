# Go 代码格式化

## 1.功能说明
* 格式化import部分：分3段式，依次为 `标准库`、`第三方库`、`项目自身库`
* 格式化单行注释：若为 `//注释内容`，调整为 `//{空格}注释内容`
* 格式化多行注释：去除首未空行；除了首未2行外，每行格式为 `{空格}*{空格}注释内容`
* 默认只对git项目库里有修改的进行格式化

对于import部分特别说明：
> 1.支持注释，注释间的空行会去除  
> 2.对于非末行的注释的位置会和其下面紧挨的importPath绑定在一起。  
> 3.末行的注释则会放入import的头部。  
> 4.对于注释的importPath,会按照其实际路径参与分组和排序。  

## 2.安装/更新
使用 go >= 1.13 版本(不需要Fan强)：
```
export GO111MODULE=on
go env GOPROXY=https://goproxy.cn,direct

go get -u github.com/fsgo/go_fmt@master
```

go < 1.13 版本(可能需要Fan强):
```
go get -u github.com/fsgo/go_fmt
```

当前版本：v0.1 20200102

## 3.使用

### 3.1 格式化git项目里有修改的.go文件
```
$ go_fmt
```

### 3.2 对当前目录所有.go文件格式化
```
$ go_fmt ./...
```

### 3.3 对指定.go文件格式化
```
$ go_fmt abc.go
```

## 4.配置到git hooks(commit前自动格式化代码)

### 4.1 配置项目Hooks
编辑项目的 `.git/hooks/pre-commit`文件，将`go_fmt`命令加入。

方式1：
```
echo -e '\ngo_fmt\n' >> $(git rev-parse --git-dir)/hooks/pre-commit
chmod 777 $(git rev-parse --git-dir)/hooks/pre-commit
```

方式2：
```
wget https://raw.githubusercontent.com/fsgo/go_fmt/master/pre-commit -O $(git rev-parse --git-dir)/hooks/pre-commit
chmod 777 $(git rev-parse --git-dir)/hooks/pre-commit
```

### 4.2 配置到全局Hooks
> 该方式会导致项目自身的hooks失效。  
> 若项目有自己的hooks，请不要配置全局而要配置到单个项目。
```
mkdir -p ~/.git_config/hooks/
git config --global core.hooksPath ~/.git_config/hooks/

wget https://raw.githubusercontent.com/fsgo/go_fmt/master/pre-commit -O ~/.git_config/hooks/pre-commit
chmod 777 ~/.git_config/hooks/pre-commit
```
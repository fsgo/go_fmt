# 自动解析当前项目的 module

## 1.推断逻辑
1.若项目包含`go.mod`文件，则读取该文件的module值(推断规则1)。
2.通过文件绝对路径来自动推断当前项目的模块名,详细如下(推断规则2)。

### 1.1 推断规则1：
当前go环境开启 go mod功能。项目存在有效的go.mod文件。   
在项目里通过如下命令可以看到项目的go.mod文件路径：
```
go env GOMOD
```

### 1.2 推断规则2：
 项目的workspace 满足如下 https://golang.org/doc/code.html#Workspaces
```
bin/
     go_fmt                               # command executable
 src/
    github.com/fsgo/go_fmt/
                           .git/          # Git repository metadata
                            go_fmt.go     # command source
 基本原理为：
 1.获取文件的完整路径
 2.查找src目录，其后第一个目录为域名
 3.默认域名再后面2级目录则为项目的模块名，如github上所有项目。
 4.若比较特殊，不是2级目录，也可以通过配置文件（~/.go_fmt/local_module.json）来设置
```

`local_module.json` 文件格式：
```
{
   "DomainLevel":{
    "abc.com":3,
   }
}
```
也就意味`abc.com` 这个域名下的项目module 规则为 `abc.com/xxx/yyy/projectzzz`
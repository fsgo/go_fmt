# go_fmt
golang format tools


## 1.前置依赖
### goimports
```
go get -u golang.org/x/tools/cmd/goimports
```

## 2.命令说明
### 2.1 bin/go_fmt.sh
使用 gofmt 和 goimports 工具对go代码进行格式化。  
import 部分会分为3部分，当前项目的会作为第3部分，顺序分别为 标准库、第三方库、项目自身库。  
import 示例：  
```go
import(
    "os"
    "log"
     
    "github.com/hidu/xxx/yyy"

    "youdomain.com/namespace/project/a"
)
```

#### useage:
```
   go_fmt.sh        # 格式化当前工作目录下，有修改的所有文件(git管理的项目)
   go_fmt.sh  all   # 格式化当前工作目录下，所有的go代码文件
   go_fmt.sh  a.go  # 格式化指定文件
```

### 2.2 bin/go_imports.sh
import 部分会分位3部分，当前项目的会作为第3部分  
####  useage:
```
  go_imports.sh a.go  # 格式化指定文件
```
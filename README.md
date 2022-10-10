# Go 代码格式化

## 1.功能说明
* 格式化 import 部分：分3段式，依默认为 `标准库`、`第三方库`、`项目自身库`
* 格式化单行注释：若为 `//注释内容`，调整为 `//{空格}注释内容`
* 默认只对 git 项目库里有修改的进行格式化
* 支持将多行的 copyright 注释修改为单行格式(默认不调整)
* 简化代码
* struct 赋值表达式，自动补齐 key
* 补充空行、移除多余的空行

对于 import 部分：
> 1.可使用`-mi`参数来控制是否将多段import合并为一段（默认否）。  
> 2.对于注释的import path,会按照其实际路径参与分组和排序。   
> 3.对于非末行的注释的位置会和其下面紧挨的import path绑定在一起。  
> 4.末行的注释则会放入import的头部。  
> 5.import path 不要使用相对路径(如`./` 和 `../`)。

会忽略当前目录以及子目录下的 `testdata` 和 `vendor` 目录。  
若需要可进入其目录里执行该命令。  

<details><summary><i>Example 1：补齐 struct key</i></summary>

``` diff
- u2 := User{"hello", 12}
+ u2 := User{Name: "hello", Age: 12}
```
</details>

<details><summary><i>Example 2：注释格式化</i></summary>

```diff
- //User 注释内容
+ // User 注释内容
type User struct{
```
</details>

<details><summary><i>Example 3：简化代码</i></summary>

1. 简化循环逻辑：
```diff
- s[a:len(s)]
+ s[a:]

- for x, _ = range v {...}
+ for x = range v {...}

- for _ = range v {...}
+ for range v {...}
```

2. 简化判断逻辑：
```diff
- if b == true {
+ if b { 

- if b == false {
+ if !b {

- if b != true {
+ if !b {

- if b != false {
+ if b {

- for b == true {
+ for b {
```

3. 使用 `strings.Contains` 替换 `strings.Count`：
```diff
- strings.Count(s, "a") == 0
+ !strings.Contains(s, "a")

- strings.Count(s, "a") > 0
+ strings.Contains(s, "a")

- strings.Count(s, "a") != 0
+ strings.Contains(s, "a")
```
`bytes.Count` 具有和 `strings.Count` 一样的规则。

4. 使用 `strings.Contains` 替换 `strings.Index`：
```diff
- strings.Index(s, "a") == -1
+ !strings.Contains(s, "a")

- strings.Index(s, "a") != -1
+ strings.Contains(s, "a")
```
`bytes.Index` 具有和 `strings.Index` 一样的规则。

5. 使用 `bytes.Equal` 替换 `bytes.Compare`：
```diff
- bytes.Compare(s,a) == 0
+ bytes.Equal(s, a)

- bytes.Compare(s,a) != 0
+ !bytes.Equal(s, a)
```

6. 使用 `==` 替换 `strings.Compare`：
```diff
- strings.Compare("abc","a") == 0
+ "abc" == "a"

- strings.Compare("abc","a") != 0
+ "abc" != "a"
```

7. 递增 1、递减 1:
```diff
- i += 1
+ i++

- i -= 1
+ i--
```

8. `time.Since`  替换 `time.Now().Sub`:
```diff
- time.Now().Sub( t1 )
+ time.Since( t1 )
```

9. `time.Until`  替换 `t.Sub( time.Now() )`:
```diff
- t1.Sub( time.Now() )
+ time.Until( t1 )
```

</details>

<details><summary><i>Example 4：重写代码</i></summary>
目前默认未启用，需使用 `go_fmt -rr` 以使用默认内置规则生效。

1. 替换废弃的 `ioutil` 的函数调用：

```diff
import (
-	"io/ioutil"
+	"io
)

- buf, err := ioutil.ReadAll(f)
+ buf, err := io.ReadAll(f)
```
</details>

<details><summary><i>Example 5：移除多余的空行</i></summary>

1. 移除 struct 内部前后多余的空行：
```diff
type userfn91 struct{
-				
	name string
-				
}
```

2. 移除 func 内部前后多余的空行：
``` diff
fn1() {
-				
	println("hello")
-				
}

```

3. 空 func 变为一行：
```diff
- fn1() {
- }
+ fn1() {}
```
</details>

<details><summary><i>Example 6：补充空行</i></summary>
在适当的位置添加空行可以增加代码的可读性。

1. struct 有文档的字段前后添加换行：

```diff
type User1 struct {
-				
	// on Name
	Name string
+				
	// on Age
	Age int
+				
	Grade int
	Class int

	Address string // 前面有空行，会保持
}
```

2. interface 有文档的方法前后添加换行：

```diff
type Group1 interface {
-				
	// Register 注册延迟函数
	Register(fn func())
+				
	Add()
+				
	// on Delete
	Delete()
+				
	Fn1()

	Fn2() // 前面有空行，会保持
-				
}
```

3. 多个定义之间添加空行:

```diff
type (
	User1 struct {
		name string
	}
+				
	User1 struct {
-				
		name string
	}
)
```

4. 全局的，不同类型定义之间添加空行：

```diff
var a="hello"
var b="world" // after b
+				
const c01="say"
+ 				
var a0 = "a0"
```

</details>



## 2.安装/更新
```bash
export GO111MODULE=on
go env GOPROXY=https://goproxy.cn,direct

go install github.com/fsgo/go_fmt@master
```
升级 Go 版本后，请用最新版本 go 重新安装/更新 `go_fmt` 。  
最低 Go 版本：go1.19


## 3.使用

### 3.0 help
> go_fmt -help

```bash
usage: go_fmt [flags] [path ...]
  -d	display diffs instead of rewriting files
  -df string
    	display diffs format, support: text, json (default "text")
  -e	enable extra rules (default true)
  -ig string
    	import group sort rule,
    	stc: Go Standard pkgs, Third Party pkgs, Current ModuleByFile pkg
    	sct: Go Standard pkgs, Current ModuleByFile pkg, Third Party pkgs
    	 (default "stc")
  -local string
    	put imports beginning with this string as 3rd-party packages (default "auto")
  -mi
    	merge imports into one
  -r value
    	rewrite rule (e.g., 'a[b:len(a)] -> a[b:]')
  -rr
    	rewrite with build in rules:
    	a[b:len(a)] -> a[b:]
    	interface{} -> any
    	a == ""     -> len(a) == 0
    	a != ""     -> len(a) != 0
    	io/#ioutil.NopCloser -> io.NopCloser
    	io/#ioutil.ReadAll   -> io.ReadAll
    	io/#ioutil.ReadFile  -> os.ReadFile
    	io/#ioutil.TempFile  -> os.CreateTemp
    	io/#ioutil.TempDir   -> os.MkdirTemp
    	io/#ioutil.WriteFile -> os.WriteFile
    	io/#ioutil.Discard   -> io.Discard
  -s	simplify code (default true)
  -slcr
    	multiline copyright to single-line
  -trace
    	show trace infos
  -w	write result to (source) file instead of stdout (default true)
```
### 3.1 格式化 `git` 项目里有修改的`.go`文件
```bash
$ go_fmt
```

### 3.2 对当前目录所有 `.go` 文件格式化
```bash
$ go_fmt ./...
```

### 3.3 对指定 `.go` 文件格式化
```bash
$ go_fmt abc.go
```

## 4.配置到 `git hooks`(commit 前自动格式化代码)

### 4.1 配置项目 Hooks
编辑项目的 `.git/hooks/pre-commit`文件，将`go_fmt`命令加入。

方式1：
```bash
echo -e '\ngo_fmt\n' >> $(git rev-parse --git-dir)/hooks/pre-commit
chmod 777 $(git rev-parse --git-dir)/hooks/pre-commit
```

方式2：
```bash
wget https://raw.githubusercontent.com/fsgo/go_fmt/master/pre-commit -O $(git rev-parse --git-dir)/hooks/pre-commit
chmod 777 $(git rev-parse --git-dir)/hooks/pre-commit
```

### 4.2 配置到全局 Hooks
> 该方式会导致项目自身的 hooks 失效。  
> 若项目有自己的 hooks，请不要配置全局而要配置到单个项目。
```bash
mkdir -p ~/.git_config/hooks/
git config --global core.hooksPath ~/.git_config/hooks/

wget https://raw.githubusercontent.com/fsgo/go_fmt/master/pre-commit -O ~/.git_config/hooks/pre-commit
chmod 777 ~/.git_config/hooks/pre-commit
```

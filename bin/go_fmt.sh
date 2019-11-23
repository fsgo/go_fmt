#!/bin/bash
set -e

###############################################################################
#  使用 gofmt 和 goimports 工具对go代码进行格式化
#  import 部分会分位3部分，当前项目的会作为第3部分
#
#   author: github.com/fsgo/go_fmt
#   since: 2019年11月23日
#   useage:
#      go_fmt.sh        # 格式化当前工作目录下，有修改的所有文件(git管理的项目)
#      go_fmt.sh  all   # 格式化当前工作目录下，所有的go代码文件
#      go_fmt.sh  a.go  # 格式化指定文件
###############################################################################


SOURCE_FILE="$1"
mkdir -p /tmp/go_fmt_sh/
find /tmp/go_fmt_sh/ -type f -name '*.go.*' -mmin 5|xargs rm -f

if [ -z `which gofmt` ];then
    echo "gofmt: missing in `PATH`"
    exit 1
fi
if [ -z `which goimports` ];then
    echo "goimports: missing in `PATH`"
    exit 1
fi

function formatImport() {
    realPath=`readlink -f "$1"`
    inFile="$2"
    
    # 获取当前文件所在package,域名后面取3层目录
    PKG=`echo "$realPath"|grep -oP '(?<=src\/)(\w+(\.\w+)+\/\w+\/\w+\/\w+)'`
    # PKG 取值，如  github.com/fsgo/abc/def/
    
    if [ -z "$PKG" ];then
        # 再次尝试域名后面取2层目录
        PKG=`echo "$realPath"|grep -oP '(?<=src\/)(\w+(\.\w+)+\/\w+\/\w+\/\w+)'`
        if [ -z "$PKG" ];then
            echo "PKG is empty" >&2
            return
        fi
    fi
    
    goimports -local "$PKG" "$inFile"
}

function goformat() {
    fileName="$1"
    
    # step 1： fix imports
    tmpFileImports="/tmp/go_fmt_sh/`basename "$fileName"`.goimports"
    rm -rf "$tmpFileImports"
    formatImport "$fileName" "$fileName" > "$tmpFileImports"
    
    
    # step 2: format code
    tmpFileFmt="/tmp/go_fmt_sh/`basename "$fileName"`.gofmt"
    rm -rf "$tmpFileFmt"
    gofmt -s=true "$tmpFileImports" > "$tmpFileFmt"
    
    if [ ! -f "$tmpFileFmt" ];then
        echo "$fileName ignore: tmpFile($tmpFileFmt) is missing"
        return
    fi
  
    if [ ! -s "$tmpFileFmt" ];then
        echo "$fileName ignore: tmpFile($tmpFileFmt) is empty"
        return
    fi
    
    lineChanges=`diff "$fileName" "$tmpFileFmt"|wc -l`
    if [ "$lineChanges" -eq "0" ];then
        printf "change_lines= %-3d\033[0m\t\033[32m%s\n" "$lineChanges" "$fileName"
    else
        cat "$tmpFileFmt" > "$fileName"
        printf "change_lines= %-3d\033[0m\t\033[34m%s\033[0m \033[5;34m\n"  "$lineChanges" "$fileName"
    fi
    printf "\033[0m"
}

# 格式化当前git项目里,修改过的文件
function formatGitWorkingDir(){
    echo "-----------------go format--------------------"
    for fName in `git status -s|grep -E "\.go$"|awk '{if($1!="D"){print $2}}'`
    do
        goformat "$fName"
    done
    echo -e "----------------------------------------------\n"
}

# 默认，无参数值，format 当前项目所有修改中的文件
if [ -z "$SOURCE_FILE" ];then
   formatGitWorkingDir
elif [ "$SOURCE_FILE" == "all" ]; then
  for fName in `find ./ -name "*.go"`
  do
      goformat "$fName"
  done
else
    goformat "$SOURCE_FILE"
fi
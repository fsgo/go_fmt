#!/bin/bash

###############################################################################
#   使用goimports 工具对go代码进行格式化
#   import 部分会分位3部分，当前项目的会作为第3部分
#
#   author: github.com/hidu
#   since: 2019年11月21日
#   useage:
#      goimports.sh code_path.go
###############################################################################

if [ ! -f "$1" ];then
    echo "file required"
    exit 1
fi

function formatImport() {
    FullName=`readlink -f  $1`
    
    echo "$FullName"
    
    if [ "${FullName##*.}"x != "go"x ];then
       echo "not go file"
       return
    fi
    
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
    
    set -x
    
    goimports -w -local "$PKG" "$FullName"
}

formatImport "$1"

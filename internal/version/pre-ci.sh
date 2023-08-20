#!/bin/bash

DAY=`date "+%Y-%m-%d"`

sed -i 's/versionDate = ".*"/versionDate = "'"${DAY}"'"/' version.go

#!/bin/bash
cd $(dirname $0)

last_sha=`git log --pretty=format:"%h"|head -n 1`
sed -i "s/github.com\/fsgo\/go_fmt@.*/github.com\/fsgo\/go_fmt@${last_sha}/g"  ../README.md
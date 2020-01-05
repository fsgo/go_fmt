#!/bin/bash
cd $(dirname $0)

last_sha=`git rev-parse --short HEAD`
sed -i "s/github.com\/fsgo\/go_fmt@.*/github.com\/fsgo\/go_fmt@${last_sha}/g"  ../README.md
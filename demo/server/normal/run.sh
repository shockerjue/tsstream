#!/bin/sh

ROOT=${PWD}
cp ${ROOT}/conf/app.conf ${ROOT}/bin/conf/
echo "${ROOT}/conf/app.conf ${ROOT}/bin/conf/"

cd ${ROOT}/bin
echo "${ROOT}/bin"

nohup ./normal > /var/log/normal.log 2>&1 &
echo "nohup ./normal > /var/log/normal.log 2>&1 &"

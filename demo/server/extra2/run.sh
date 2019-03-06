#!/bin/sh

ROOT=${PWD}
cp ${ROOT}/conf/app.conf ${ROOT}/bin/conf/
echo "${ROOT}/conf/app.conf ${ROOT}/bin/conf/"

cd ${ROOT}/bin
echo "${ROOT}/bin"

nohup ./extra2 > /var/log/extra2.log 2>&1 &
echo "nohup ./extra2 > /var/log/extra2.log 2>&1 &"

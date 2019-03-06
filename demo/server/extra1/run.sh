#!/bin/sh

ROOT=${PWD}
cp ${ROOT}/conf/app.conf ${ROOT}/bin/conf/
echo "${ROOT}/conf/app.conf ${ROOT}/bin/conf/"

cd ${ROOT}/bin
echo "${ROOT}/bin"

nohup ./extra1 > /var/log/extra1.log 2>&1 &
echo "nohup ./extra1 > /var/log/extra1.log 2>&1 &"

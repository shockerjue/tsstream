#!/bin/sh

ROOT=${PWD}
cp ${ROOT}/conf/app.conf ${ROOT}/bin/conf/
echo "${ROOT}/conf/app.conf ${ROOT}/bin/conf/"

cd ${ROOT}/bin
echo "${ROOT}/bin"

nohup ./monitor > /var/log/monitor.log 2>&1 &
echo "nohup ./monitor > /var/log/monitor.log 2>&1 &"
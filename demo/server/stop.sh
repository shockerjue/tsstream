#!/bin/sh

ROOT_PWD=$(pwd)
echo ${ROOT_PWD}

cd ${ROOT_PWD}/extra1
sh stop.sh
sleep 5

cd ${ROOT_PWD}/extra2
sh stop.sh
sleep 5

cd ${ROOT_PWD}/normal
sh stop.sh

#!/bin/sh

ps -ef | grep "./extra1" | grep -v grep | awk '{print $2}' | xargs kill -9
echo "ps -ef | grep "./extra1" | grep -v grep | awk '{print $2}' | xargs kill -9"

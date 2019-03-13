#!/bin/sh

ps -ef | grep "./monitor" | grep -v grep | awk '{print $2}' | xargs kill -9
echo "ps -ef | grep "./monitor" | grep -v grep | awk '{print $2}' | xargs kill -9"


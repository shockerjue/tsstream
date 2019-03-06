#!/bin/sh

ps -ef | grep "./normal" | grep -v grep | awk '{print $2}' | xargs kill -9
echo "ps -ef | grep "./normal" | grep -v grep | awk '{print $2}' | xargs kill -9"

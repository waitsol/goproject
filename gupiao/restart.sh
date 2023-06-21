#!/bin/bash

# 查找进程ID
pid=$(pgrep gupiao)

if [ -z "$pid" ]; then
  echo "进程 'gupiao' 不存在"
else
  # 杀死进程
  kill "$pid"
  echo "进程 'gupiao' (PID: $pid) 已被杀死"
fi
rm *.out
go build -o gupiao
if [ "$#" -eq 0 ]; then
  demon ./gupiao
  echo "启动服务"
else
  echo "没有启动服务"
fi
#!/bin/bash

POD_IP=$(curl ifconfig.me)  # 获取容器的 IP

# 将 POD_IP 设置为环境变量
export POD_IP

# 启动主应用
exec "$RUN_BIN -f $RUN_CONF"
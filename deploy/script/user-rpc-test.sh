#!/bin/bash
reso_addr='registry.cn-hangzhou.aliyuncs.com/easy-chat-xy/user-rpc-dev'
tag='latest'



container_name="easy-chat-user-rpc-test"

POD_IP="8.210.39.173"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-chat -v /easy-chat/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10000:10000 -e pod_ip=${POD_IP}  -v /root/easy-chat/user-rpc-logs:/user/logs/ --name=${container_name} -d ${reso_addr}:${tag}

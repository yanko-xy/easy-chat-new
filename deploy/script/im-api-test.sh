#!/bin/bash

server_name="im"
server_type="api"
port = 8890

reso_addr="registry.cn-hangzhou.aliyuncs.com/easy-chat-xy/${server_name}-${server_type}-dev"
tag='latest'


container_name="easy-chat-${server_name}-${server_type}-test"

pod_ip="120.26.209.19"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-chat -v /easy-chat/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p ${port}:${port} -v /root/easy-chat/${server_name}-${server_type}-logs:/${server_name}/logs/ --name=${container_name} -d ${reso_addr}:${tag}

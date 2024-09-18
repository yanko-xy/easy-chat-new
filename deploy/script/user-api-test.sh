#!/bin/bash
reso_addr='registry.cn-hangzhou.aliyuncs.com/easy-chat-xy/user-api-dev'
tag='latest'

container_name="easy-chat-user-api-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-chat -v /easy-chat/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 8888:8888 -v /root/easy-chat/user-api-logs:/user/logs/ --name=${container_name} -d ${reso_addr}:${tag}

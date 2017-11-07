docker plugin disable redislog:latest 
docker plugin rm redislog:latest 
docker plugin create redislog /root/go/src/h3d.com/weipeng/dockerlogredis
docker plugin enable redislog

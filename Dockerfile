FROM  golang:1.7

COPY . /go/src/h3d.com/weipeng/dockerlogredis
RUN cd /go/src/h3d.com/weipeng/dockerlogredis && go get && go build --ldflags '-extldflags "-static"' -o /usr/bin/dockerlogredis

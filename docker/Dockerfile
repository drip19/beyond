# FROM...as...和下面的COPY --from连用
FROM golang:1.15.2-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

#COPY 是把本地的文件拷贝到容器镜像中
#COPY <src> <dest>
# 注意这个源路径是项目跟路径
COPY . /chandler/beyond
RUN cd /chandler/beyond && make clean build


FROM alpine:3.12.0

ENV Chandler /chandler
RUN apk --no-cache add ca-certificates
WORKDIR $Chandler
COPY --from=builder /chandler/beyond/build/gbeyond /usr/local/bin/gbeyond

#HEALTHCHECK --interval=30s --timeout=30s --retries=3 CMD curl -f http://127.0.0.1:29705/ping/info

# ENTRYPOINT： container启动时执行的命令，而且一定会被执行，如果一个Dockerfile文件中有多个，只有最后一个生效
# 这里就是启动编译出来的gmydemo文件
ENTRYPOINT [ "gbeyond"]


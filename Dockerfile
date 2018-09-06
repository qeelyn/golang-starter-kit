FROM alpine:latest

ENV TIMEZONE=Asia/Shanghai HOSTIP=0.0.0.0
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add ca-certificates tzdata && \
    echo "$TIMEZONE" > /etc/timezone && \
    ln -sf "/usr/share/zoneinfo/$TIME_ZONE" /etc/localtime
# in alpine,cross compile binany file maynot appear 'sh: file not found ' error
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

VOLUME /app/runtime

WORKDIR /app

COPY cmd/serve .
COPY cmd/config ./config
COPY cmd/public ./public
CMD ./serve all
FROM golang:1.17-buster as builder


LABEL name="jd_cookie server"
LABEL version="3.0.1"
LABEL author="scjtqs <scjtqs@qq.com>"
LABEL maintainer="scjtqs <scjtqs@qq.com>"
LABEL description="simple to get jd cookie"
COPY ./sources.list /etc/apt/sources.list
ARG Version="v3.0.1"

ADD . /src
#ENV GOPROXY "http://goproxy.cn,direct"
ENV CGO_ENABLED "0"
ENV GO111MODULE "on"



##替换官方源为国内源
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN cd /src \
    && apt-get update \
    && apt-get install -y build-essential openssl git\
    && rm -rf /var/lib/apt/lists/*apt \
    && rm -rf dist \
    && go mod tidy \
    && make

FROM alpine:3.13 as production

ENV UPSAVE ""
ENV UPSAVE_METHOD "POST" # POST GET
ENV UPSAVE_KEY  "userCookie" # post OR  get 方式的 cookie传递的key

ENV DB_ENABLE "false"
ENV DB_HOST ""
ENV DB_PORT ""
ENV DB_USER ""
ENV DB_PASS ""
ENV DB_DATABASE ""
ENV DB_TYPE "mysql"

COPY --from=builder /src/dist /opt/app

RUN  adduser -D -H www \
     && chown -R www /opt/app \
     && apk add -U --no-cache tzdata \
     && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
     && apk del tzdata

USER www
WORKDIR /opt/app

EXPOSE 29099

CMD ["/opt/app/jdcookie"]
# Import sql manually https://github.com/zhenorzz/goploy/blob/master/model/sql/goploy.sql
FROM alpine
LABEL maintainer="zhenorzz@gmail.com"
ARG GOPLOY_VER=v1.9.0
ENV GOPLOY_VER=${GOPLOY_VER}

ENV MYSQL_PORT=3306

RUN echo "http://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/latest-stable/community/" >> /etc/apk/repositories

#install
RUN apk update && \
    apk add --no-cache \
    openssh-client \
    ca-certificates \
    bash \
    git \
    rsync \
    && rm -rf /var/cache/apk/* 

#git
RUN git config --global pull.rebase false

#goploy
ADD https://github.com/zhenorzz/goploy/releases/download/${GOPLOY_VER}/goploy /opt/goploy/
RUN chmod a+x /opt/goploy/goploy

EXPOSE 80

VOLUME ["/opt/goploy/repository/"]

WORKDIR /opt/goploy/

ENTRYPOINT ["bash", "-c"]

CMD ["./goploy --asset-dir=./repository"]

# Import sql manually https://github.com/zhenorzz/goploy/blob/master/model/sql/goploy.sql
FROM alpine
LABEL maintainer="zhenorzz@gmail.com"
ARG GOPLOY_VER=v1.16.2
ENV GOPLOY_VER=${GOPLOY_VER}

RUN echo "https://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories
RUN echo "https://mirrors.aliyun.com/alpine/latest-stable/community/" >> /etc/apk/repositories

#install
RUN apk update && \
    apk add --no-cache \
    openssh-client \
    ca-certificates \
    sshpass \
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

ENTRYPOINT /opt/goploy/goploy --asset-dir=/opt/goploy/repository

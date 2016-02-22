FROM ubuntu:14.04
MAINTAINER l00374667 "l00273667@openvmse.org"
ENV REFRESHED_AT 2016-02-17

RUN mkdir -p /opt/service

COPY serv-registrator /opt/service/serv-registrator
COPY config.ini /opt/service/config.ini

WORKDIR /opt/service

ENTRYPOINT ["./serv-registrator"]

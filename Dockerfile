FROM ubuntu:14.04
MAINTAINER l00374667 "l00273667@openvmse.org"
ENV REFRESHED_AT 2016-02-17

RUN mkdir -p /opt/service

COPY serv_reg /opt/service/serv_reg

WORKDIR /opt/service

ENTRYPOINT ["./serv_reg"]

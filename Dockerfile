#
FROM ubuntu:18.04 

RUN mkdir -p /opt/playground /data

WORKDIR /opt/playground

COPY ./ingester/ingester ingester

EXPOSE 9080
VOLUME /data

ENTRYPOINT ["/opt/playground/ingester"]


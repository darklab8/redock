FROM ubuntu:20.04

ARG BUILD_VERSION
ENV BUILD_VERSION=${BUILD_VERSION}

CMD sleep infinity

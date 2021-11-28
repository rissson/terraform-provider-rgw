FROM golang:1.16

ENV DEBIAN_FRONTEND noninteractive
ENV LC_ALL C.UTF-8
ENV LANG C.UTF-8

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
            curl ca-certificates gnupg apt-transport-https git software-properties-common

RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc make

RUN curl -1sLf https://download.ceph.com/keys/release.asc | apt-key add - && \
    echo "deb https://download.ceph.com/debian-octopus/ bullseye main" > /etc/apt/sources.list.d/ceph.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends libcephfs-dev librbd-dev librados-dev

RUN go get -u golang.org/x/lint/golint
RUN go get -u golang.org/x/tools/...

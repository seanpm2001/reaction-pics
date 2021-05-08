FROM ubuntu:20.04

LABEL maintainer="git@albertyw.com"
EXPOSE 5003
HEALTHCHECK --interval=5s --timeout=3s CMD bin/healthcheck.sh || exit 1
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Set locale
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV DEBIAN_FRONTEND noninteractive

# Install go and other dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential locales                      `: basic packages` \
    git curl tar ca-certificates                 `: go installation` \
    gcc g++ make gnupg                           `: nodejs dependencies` \
    && curl https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz | tar xvz \
    && ./godeb install "$(./godeb list | tail -n 1)" \
    && rm godeb \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 \
    && curl https://deb.nodesource.com/setup_16.x | bash \
    && apt-get install -y --no-install-recommends \
    nodejs                                      ` Javascript assets` \
    && rm -rf /var/lib/apt/lists/*

# Set up directory structures
ENV GOPATH /root/gocode
RUN mkdir -p /root/gocode/src/github.com/albertyw/reaction-pics
COPY . /root/gocode/src/github.com/albertyw/reaction-pics
WORKDIR /root/gocode/src/github.com/albertyw/reaction-pics

# App-specific setup
RUN bin/container_setup.sh

CMD ["bin/start.sh"]

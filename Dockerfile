FROM ubuntu:20.04

LABEL maintainer="git@albertyw.com"
EXPOSE 5003
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Set locale
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV DEBIAN_FRONTEND noninteractive

# Install go and other dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential locales                      `: basic packages` \
    git curl wget tar ca-certificates            `: go installation` \
    gcc g++ make gnupg                           `: nodejs dependencies` \
    && wget -nv https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz \
    && tar xvf godeb-amd64.tar.gz \
    && ./godeb install "$(./godeb list | tail -n 1)" \
    && rm godeb* go_*deb \
    && sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && locale-gen \
    && wget -nv https://deb.nodesource.com/setup_16.x && bash setup_16.x \
    && apt-get install -y --no-install-recommends nodejs && \
    rm -rf /var/lib/apt/lists/*

# Set up directory structures
ENV GOPATH /root/gocode
RUN mkdir -p /root/gocode/src/github.com/albertyw/reaction-pics
COPY . /root/gocode/src/github.com/albertyw/reaction-pics
WORKDIR /root/gocode/src/github.com/albertyw/reaction-pics

# App-specific setup
RUN bin/container_setup.sh

CMD ["bin/start.sh"]

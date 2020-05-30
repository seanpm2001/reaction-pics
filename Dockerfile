# Needed to match supported versions of golang-backports
# TODO - upgrade this to 20.04 once https://launchpad.net/~longsleep/+archive/ubuntu/golang-backports supports it
FROM ubuntu:18.04
LABEL maintainer="git@albertyw.com"
EXPOSE 5003

# Set locale
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV DEBIAN_FRONTEND noninteractive

# Install go and other dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential locales        `: basic packages` \
    git curl tar ca-certificates   `: go installation` \
    gcc g++ make gnupg             `: nodejs dependencies`
RUN curl https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz -o godeb-amd64.tar.gz && \
    tar xvf godeb-amd64.tar.gz &&\
    ./godeb install "$(./godeb list | tail -n 1)"
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && locale-gen

# Install node
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get update && apt-get install -y --no-install-recommends nodejs
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Set up directory structures
ENV GOPATH /root/gocode
RUN mkdir -p /root/gocode/src/github.com/albertyw/reaction-pics
COPY . /root/gocode/src/github.com/albertyw/reaction-pics
WORKDIR /root/gocode/src/github.com/albertyw/reaction-pics

# App-specific setup
RUN bin/container_setup.sh

CMD ["bin/start.sh"]

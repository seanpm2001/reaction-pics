FROM debian:buster
LABEL maintainer="git@albertyw.com"
EXPOSE 5003

# Set locale
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

# Install updates and system packages
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update && RUN apt-get install -y --no-install-recommends \
    build-essential \
    locales \
    software-properties-common \
    curl \
    golang-go \
    git

# Install node
RUN apt-get install -y --no-install-recommends gcc g++ make
RUN curl -sL https://deb.nodesource.com/setup_11.x | bash -
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

FROM node:18 as node
WORKDIR /
COPY . .
RUN npm ci --only=production \
    && npm run minify \
    && sed -i '' server/static/**/*


FROM golang:1.18-bullseye

LABEL maintainer="git@albertyw.com"
EXPOSE 5003
HEALTHCHECK --interval=5s --timeout=3s CMD bin/healthcheck.sh || exit 1
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Set locale
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
    locales                                     `: Basic-packages` \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 \
    && rm -rf /var/lib/apt/lists/*

# Set up directory structures
WORKDIR /root/gocode/src/github.com/albertyw/reaction-pics
ENV GOPATH /root/gocode
RUN mkdir -p .
COPY . .
COPY --from=node ./server/static ./server/static

# App-specific setup
RUN make bins

CMD ["bin/start.sh"]

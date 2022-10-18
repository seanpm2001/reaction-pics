FROM node:18 as node
WORKDIR /
COPY . .
RUN npm ci --only=production \
    && npm run minify \
    && sed -i '' server/static/**/*


FROM golang:1.19-bullseye as go
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Set up directory structures
WORKDIR /root/gocode/src/github.com/albertyw/reaction-pics
ENV GOPATH /root/gocode
RUN mkdir -p .
COPY . .
COPY --from=node ./server/static ./server/static

# App-specific setup
RUN make bins

FROM alpine:3
LABEL maintainer="git@albertyw.com"
EXPOSE 5003
HEALTHCHECK --interval=5s --timeout=3s CMD bin/healthcheck.sh || exit 1

WORKDIR /root/
RUN apk add --no-cache libc6-compat=1.2.3-r0
COPY --from=go /root/gocode/src/github.com/albertyw/reaction-pics/reaction-pics .
COPY --from=go /root/gocode/src/github.com/albertyw/reaction-pics/.env* ./
RUN mkdir -p /root/gocode/src/github.com/albertyw/reaction-pics/logs/app

CMD ["./reaction-pics"]

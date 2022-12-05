FROM node:18 as node
WORKDIR /
COPY . .
RUN npm ci --only=production \
    && npm run minify \
    && sed -i '' server/static/**/*


FROM golang:1.19-bullseye as go
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Set up directory structures
WORKDIR /root/
RUN mkdir -p .
COPY . .
COPY --from=node ./server/static ./server/static

# App-specific setup
RUN make bins

FROM alpine:3
LABEL maintainer="git@albertyw.com"
EXPOSE 5003
HEALTHCHECK --interval=5s --timeout=3s CMD ./healthcheck.sh || exit 1

WORKDIR /root/
RUN apk add --no-cache libc6-compat=1.2.3-r4 curl==7.86.0-r1
COPY --from=go /root/reaction-pics .
COPY --from=go /root/.env* ./
COPY --from=go /root/bin/healthcheck.sh .
RUN mkdir -p /root/logs/app

CMD ["./reaction-pics"]

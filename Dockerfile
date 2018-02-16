FROM alpine:3.6
RUN apk update && apk add --no-cache \
    ca-certificates \
    git \
    && update-ca-certificates
COPY bin/gh-server /usr/local/bin/gh-server
CMD gh-server

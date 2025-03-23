# Based on https://github.com/danderson/netboot/blob/main/dockerfiles/pixiecore/Dockerfile
FROM alpine:3.21

RUN set -x                                                                      && \
    set -e                                                                      && \
    apk upgrade --update-cache                                                  && \
    apk add ca-certificates                                                     && \
    apk add --virtual .build-deps git go musl-dev

RUN NAMESPACE=go.universe.tf                                                    && \
    REPO=netboot                                                                && \
    PKG=cmd/pixiecore                                                           && \
    NAMESPACE_PATH="$GOPATH/src/$NAMESPACE"                                     && \
    REPO_PATH="$NAMESPACE_PATH/$REPO"                                           && \
    PKG_PATH="$REPO_PATH/$PKG"                                                  && \
    GOBIN=/usr/local/bin go install "$NAMESPACE/$REPO/$PKG@latest"
    #apk del --purge .build-deps                                                 && \
    #rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/pixiecore"]
FROM bitmark/nodejs-env as build-client

COPY ui /go/src/github.com/bitmark-inc/bitmark-node/ui
RUN cd /go/src/github.com/bitmark-inc/bitmark-node/ui && bash -c "source ~/.nvm/nvm.sh && npm install && npm run build"

# Dockerfile after bitmarkd v0.9.0 for adapting go module
#--

FROM bitmark/go-env:ubuntu1804 as go-env

# VERSION SHOW ON BITMARK-NODE
ENV VERSION v1.0
ENV BITMARKD_VERSION v0.10.0

RUN apt-get install libargon2-0-dev

RUN go get -d github.com/bitmark-inc/bitmarkd || \
    cd /go/src/github.com/bitmark-inc/bitmarkd && \
    git checkout "$BITMARKD_VERSION"

ENV GO111MODULE on

RUN cd /go/src/github.com/bitmark-inc/bitmarkd && \
    go mod download && \
    go install -ldflags "-X main.version=$BITMARKD_VERSION" github.com/bitmark-inc/bitmarkd/command/...

ENV GO111MODULE off

RUN go get -d github.com/bitmark-inc/go-argon2 && \
    go get -d github.com/bitmark-inc/go-libucl && \
    go get github.com/bitmark-inc/exitwithstatus && \
    go get github.com/bitmark-inc/bitmark-sdk-go

RUN go get github.com/bitmark-inc/discovery && \
    go get -d github.com/bitmark-inc/bitmark-wallet && \
    go install github.com/bitmark-inc/bitmark-wallet

RUN go get github.com/gin-gonic/gin && go get github.com/gin-gonic/contrib/static && \
    go get github.com/coreos/bbolt

COPY . /go/src/github.com/bitmark-inc/bitmark-node
RUN go install -ldflags "-X main.version=$VERSION" github.com/bitmark-inc/bitmark-node
COPY --from=build-client /go/src/github.com/bitmark-inc/bitmark-node/ui/public/ /go/src/github.com/bitmark-inc/bitmark-node/ui/public/

ADD bitmark-node.conf.sample /.config/bitmark-node/bitmark-node.conf
ADD docker-assets/bitmarkd.conf /.config/bitmark-node/bitmarkd/bitmark/
ADD docker-assets/recorderd.conf /.config/bitmark-node/recorderd/bitmark/
ADD docker-assets/bitmarkd-test.conf /.config/bitmark-node/bitmarkd/testing/bitmarkd.conf
ADD docker-assets/recorderd-test.conf /.config/bitmark-node/recorderd/testing/recorderd.conf
ADD docker-assets/start.sh /

ENV NETWORK bitmark

EXPOSE 2130 2131 2135 2136
CMD ["/start.sh"]

FROM bitmark/nodejs-env as build-client

COPY ui /go/src/bitmark-node-docker/ui
RUN cd /go/src/bitmark-node-docker/ui && bash -c "source ~/.nvm/nvm.sh && npm install && npm run build"

FROM bitmark/go-env:go12 as go-env

# VERSION SHOW ON BITMARK-NODE
ENV VERSION v1.3.0
ENV BITMARKD_VERSION v0.11.0-rc.2

# Install argon2 for OS
RUN apt-get install libargon2-0-dev && apt-get install -y jq

# Get Bitmarkd and corresponding version


# Install and build bitmark-cli  bitmark-dumpdb  bitmark-info  bitmarkd  recorderd
ENV GO111MODULE on

RUN cd /go/src && \
    git clone --branch="$BITMARKD_VERSION" https://github.com/bitmark-inc/bitmarkd.git && \
    git clone https://github.com/bitmark-inc/discovery && \
    git clone https://github.com/bitmark-inc/bitmark-wallet

RUN mkdir /go/src/bitmark-node-docker
COPY . /go/src/bitmark-node-docker

RUN cd /go/src/bitmarkd && \
    go install -ldflags "-X main.version=$BITMARKD_VERSION" ./command/... && \
    cd /go/src/bitmark-node-docker && \
    go install 

# COPY static ui to bitmark-node-docker
COPY --from=build-client /go/src/bitmark-node-docker/ui/public/ /go/src/bitmark-node-docker/ui/public/

ADD bitmark-node-docker.conf.sample /.config/bitmark-node-docker/bitmark-node-docker.conf
ADD docker-assets/bitmarkd.conf /.config/bitmark-node-docker/bitmarkd/bitmark/
ADD docker-assets/recorderd.conf /.config/bitmark-node-docker/recorderd/bitmark/
ADD docker-assets/bitmarkd-test.conf /.config/bitmark-node-docker/bitmarkd/testing/bitmarkd.conf
ADD docker-assets/recorderd-test.conf /.config/bitmark-node-docker/recorderd/testing/recorderd.conf
ADD docker-assets/start.sh /

ENV NETWORK bitmark

EXPOSE 2130 2131 2135 2136
CMD ["/start.sh"]

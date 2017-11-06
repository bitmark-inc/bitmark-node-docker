FROM ubuntu

RUN apt-get -q update && \
    apt-get -yq install automake autoconf pkg-config libtool software-properties-common && \
    apt-get -yq install git wget net-tools vim && \
    apt-get -y autoclean

RUN git clone https://github.com/vstakhov/libucl /libucl && \
    cd /libucl && \
    ./autogen.sh && \
    ./configure --disable-debug --disable-dependency-tracking --disable-silent-rules --prefix=/usr/local/ && \
    make install && \
    rm -rf /libucl

RUN mkdir /phc-winner-argon2 && \
    wget -qO- https://github.com/P-H-C/phc-winner-argon2/archive/20161029.tar.gz | tar zx --strip-components 1 -C /phc-winner-argon2 && \
    cd /phc-winner-argon2 && \
    make && make install PREFIX=/usr/local && \
    cp /phc-winner-argon2/*.pc /usr/lib/pkgconfig/ && \
    rm -rf /phc-winner-argon2

RUN add-apt-repository ppa:longsleep/golang-backports && apt-get -q update && apt-get install -yq golang-go libzmq3-dev

ENV GOPATH /go
ENV PATH="/go/bin:${PATH}"

RUN wget -qO- https://raw.githubusercontent.com/creationix/nvm/v0.33.2/install.sh | bash && \
    bash -c "source ~/.nvm/nvm.sh && nvm install v7"

ENV BITMARKD_VERSION 6.2

RUN go get -d github.com/bitmark-inc/bitmarkd || \
    cd /go/src/github.com/bitmark-inc/bitmarkd && git checkout v6.2 && \
    go install -ldflags "-X main.version=$BITMARKD_VERSION" github.com/bitmark-inc/bitmarkd/command/... && \
    go get github.com/bitmark-inc/discovery && \
    go get -d github.com/bitmark-inc/bitmark-wallet && \
    go install github.com/bitmark-inc/bitmark-wallet

RUN go get github.com/bitmark-inc/exitwithstatus && \
    go get github.com/gin-gonic/gin && go get github.com/gin-gonic/contrib/static && \
    go get github.com/coreos/bbolt

COPY . /go/src/github.com/bitmark-inc/bitmark-node
RUN go install github.com/bitmark-inc/bitmark-node
RUN cd /go/src/github.com/bitmark-inc/bitmark-node/ui && bash -c "source ~/.nvm/nvm.sh && npm install && npm run build"

ADD bitmark-node.conf.sample /.config/bitmark-node/bitmark-node.conf
ADD docker-assets/bitmarkd.conf /.config/bitmark-node/bitmarkd/bitmark/
ADD docker-assets/recorderd.conf /.config/bitmark-node/recorderd/bitmark/
ADD docker-assets/bitmarkd-test.conf /.config/bitmark-node/bitmarkd/testing/bitmarkd.conf
ADD docker-assets/recorderd-test.conf /.config/bitmark-node/recorderd/testing/recorderd.conf
ADD docker-assets/start.sh /

EXPOSE 2130 2135 2136 2150
CMD ["/start.sh"]

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

RUN go get -d github.com/bitmark-inc/bitmarkd || \
    go install github.com/bitmark-inc/bitmarkd/command/... && \
    go get github.com/bitmark-inc/discovery && \
    go get -d github.com/bitmark-inc/bitmark-wallet && \
    go install github.com/bitmark-inc/bitmark-wallet

ADD docker-assets/bitmarkd.conf /.config/bitmarkd/
ADD docker-assets/prooferd.conf /.config/prooferd/
ADD docker-assets/discovery.conf /.config/discovery/

RUN bitmarkd --config-file=/.config/bitmarkd/bitmarkd.conf gen-peer-identity && \
    bitmarkd --config-file=/.config/bitmarkd/bitmarkd.conf gen-rpc-cert && \
    bitmarkd --config-file /.config/bitmarkd/bitmarkd.conf gen-proof-identity && \
    prooferd --config-file /.config/prooferd/prooferd.conf generate-identity

EXPOSE 2130 2135 2136 2150
CMD ["/go/bin/bitmarkd", "--config-file", "/.config/bitmarkd/bitmarkd.conf"]

#!/bin/sh


# Set the data-directory path for bitmarkd and recorder configuration 
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/bitmarkd/bitmark:' /.config/bitmark-node-docker/bitmarkd/bitmark/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/recorderd/bitmark:' /.config/bitmark-node-docker/recorderd/bitmark/recorderd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/bitmarkd/testing:' /.config/bitmark-node-docker/bitmarkd/testing/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/recorderd/testing:' /.config/bitmark-node-docker/recorderd/testing/recorderd.conf

# Generate all keys for livenet
# @ /.config/bitmark-node-docker/bitmarkd/bitmark
cd /.config/bitmark-node-docker/bitmarkd/bitmark
bitmarkd --config-file=/.config/bitmark-node-docker/bitmarkd/bitmark/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node-docker/bitmarkd/bitmark/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node-docker/bitmarkd/bitmark/bitmarkd.conf gen-proof-identity
# copy peer keys to recorderd
cp peer.private /.config/bitmark-node-docker/recorderd/bitmark/peer.private
cp peer.public /.config/bitmark-node-docker/recorderd/bitmark/peer.public
# @ /.config/bitmark-node-docker/recorderd/bitmark
cd /.config/bitmark-node-docker/recorderd/bitmark
recorderd --config-file=/.config/bitmark-node-docker/recorderd/bitmark/recorderd.conf generate-identity
# remove proof.sign for livenet
rm /.config/bitmark-node-docker/bitmarkd/bitmark/proof.sign


# Generate all keys for testnet
# @ /.config/bitmark-node-docker/bitmarkd/testing
cd /.config/bitmark-node-docker/bitmarkd/testing/
bitmarkd --config-file=/.config/bitmark-node-docker/bitmarkd/testing/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node-docker/bitmarkd/testing/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node-docker/bitmarkd/testing/bitmarkd.conf gen-proof-identity
# copy peer keys to recorderd
cp peer.private /.config/bitmark-node-docker/recorderd/testing/peer.private
cp peer.public /.config/bitmark-node-docker/recorderd/testing/peer.public
# @ /.config/bitmark-node-docker/recorderd/testing
cd /.config/bitmark-node-docker/recorderd/testing
recorderd --config-file=/.config/bitmark-node-docker/recorderd/testing/recorderd.conf generate-identity
# remove proof.sign for testnet
rm /.config/bitmark-node-docker/bitmarkd/testing/proof.sign

# move back to root directory
cd /

# Set the proof public key inot recorderd config
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node-docker/bitmarkd/bitmark/proof.public | cut -d":" -f2)/g" /.config/bitmark-node-docker/recorderd/bitmark/recorderd.conf
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node-docker/bitmarkd/testing/proof.public | cut -d":" -f2)/g" /.config/bitmark-node-docker/recorderd/testing/recorderd.conf

# Set the data-directory path for bitmarkd and recorder configuration 
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/bitmarkd/bitmark:' /.config/bitmark-node-docker/bitmarkd/bitmark/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/bitmarkd/testing:' /.config/bitmark-node-docker/bitmarkd/testing/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/recorderd/bitmark:' /.config/bitmark-node-docker/recorderd/bitmark/recorderd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node-docker/recorderd/testing:' /.config/bitmark-node-docker/recorderd/testing/recorderd.conf

# Start bitmark node
export CONTAINER_IP=$(awk 'END{print $1}' /etc/hosts)
cd /.config/bitmark-node-docker
bitmark-node-docker -config-file=/.config/bitmark-node/bitmark-node-docker.conf -container-ip=$CONTAINER_IP -ui=/go/src/bitmark-node-docker/ui/public

#!/bin/sh

# Generate all keys
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-proof-identity
recorderd --config-file=/.config/bitmark-node/recorderd/bitmark/recorderd.conf generate-identity
rm /.config/bitmark-node/bitmarkd/bitmark/proof.sign

# Copy livenet key to its data directory
cp peer.private /.config/bitmark-node/bitmarkd/bitmark/peer.private
cp peer.public /.config/bitmark-node/bitmarkd/bitmark/peer.public
cp proof.private /.config/bitmark-node/bitmarkd/bitmark/proof.private
cp proof.public /.config/bitmark-node/bitmarkd/bitmark/proff.public
cp rpc.crt /.config/bitmark-node/bitmarkd/bitmark/rpc.crt
cp rpc.key /.config/bitmark-node/bitmarkd/bitmark/rpc.key


bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-proof-identity
recorderd --config-file=/.config/bitmark-node/recorderd/testing/recorderd.conf generate-identity
rm /.config/bitmark-node/bitmarkd/testing/proof.sign

# copy testnet key to its data directory
cp peer.private /.config/bitmark-node/bitmarkd/testing/peer.private
cp peer.public /.config/bitmark-node/bitmarkd/testing/peer.public
cp proof.private /.config/bitmark-node/bitmarkd/testing/proof.private
cp proof.public /.config/bitmark-node/bitmarkd/testing/proof.public
cp rpc.crt /.config/bitmark-node/bitmarkd/testing/rpc.crt
cp rpc.key /.config/bitmark-node/bitmarkd/testing/rpc.key

# Set the proof public key inot recorderd config
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node/bitmarkd/bitmark/proof.public | cut -d":" -f2)/g" /.config/bitmark-node/recorderd/bitmark/recorderd.conf
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node/bitmarkd/testing/proof.public | cut -d":" -f2)/g" /.config/bitmark-node/recorderd/testing/recorderd.conf

# Set the data-directory path for bitmarkd configuration 
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/bitmarkd/bitmark:' /.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/bitmarkd/testing:' /.config/bitmark-node/bitmarkd/testing/bitmarkd.conf

# Start bitmark node
export CONTAINER_IP=$(awk 'END{print $1}' /etc/hosts)
cd /.config/bitmark-node
bitmark-node -config-file=/.config/bitmark-node/bitmark-node.conf -container-ip=$CONTAINER_IP -ui=/go/src/github.com/bitmark-inc/bitmark-node/ui/public

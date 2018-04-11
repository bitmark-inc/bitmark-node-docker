#!/bin/sh

# Generate all keys
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-proof-identity
recorderd --config-file=/.config/bitmark-node/recorderd/bitmark/recorderd.conf generate-identity
rm /.config/bitmark-node/bitmarkd/bitmark/proof.sign

bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-proof-identity
recorderd --config-file=/.config/bitmark-node/recorderd/testing/recorderd.conf generate-identity
rm /.config/bitmark-node/bitmarkd/testing/proof.sign

# Set the proof public key inot recorderd config
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node/bitmarkd/bitmark/proof.public | cut -d":" -f2)/g" /.config/bitmark-node/recorderd/bitmark/recorderd.conf
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node/bitmarkd/testing/proof.public | cut -d":" -f2)/g" /.config/bitmark-node/recorderd/testing/recorderd.conf

# Start bitmark node
export CONTAINER_IP=$(awk 'END{print $1}' /etc/hosts)
cd /.config/bitmark-node
bitmark-node -config-file=/.config/bitmark-node/bitmark-node.conf -container-ip=$CONTAINER_IP -ui=/go/src/github.com/bitmark-inc/bitmark-node/ui/public

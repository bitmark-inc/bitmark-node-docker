#!/bin/sh

# Generate all keys
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-proof-identity
prooferd --config-file=/.config/bitmark-node/prooferd/bitmark/prooferd.conf generate-identity

bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-proof-identity
prooferd --config-file=/.config/bitmark-node/prooferd/testing/prooferd.conf generate-identity

# Set the proof public key inot prooferd config
sed -ie "s/@BITMARKD-PROOFER-PUBLIC-KEY@/$(cat /.config/bitmark-node/bitmarkd/proof.public | cut -d":" -f2)/g" /.config/bitmark-node/prooferd/prooferd.conf

# Start bitmark node
export CONTAINER_IP=$(awk 'END{print $1}' /etc/hosts)
bitmark-node -config-file=/.config/bitmark-node/bitmark-node.conf -container-ip=$CONTAINER_IP -ui=/go/src/github.com/bitmark-inc/bitmark-node/ui/public

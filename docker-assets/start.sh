#!/bin/sh

# Generate all keys
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmarkd.conf gen-proof-identity
prooferd --config-file=/.config/bitmark-node/prooferd/prooferd.conf generate-identity

# Set the proof public key inot prooferd config
sed -ie "s/@BITMARKD-PROOFER-PUBLIC-KEY@/$(cat /.config/bitmarkd/proof.public | cut -d":" -f2)/g" /.config/prooferd/prooferd.conf

# Start bitmark node
export CONTAINER_IP=$(awk 'END{print $1}' /etc/hosts)
bitmark-node -config-file=/.config/bitmark-node/bitmark-node.conf -container-ip=$CONTAINER_IP

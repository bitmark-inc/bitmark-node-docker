#!/bin/sh


# Set the data-directory path for bitmarkd and recorder configuration 
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/bitmarkd/bitmark:' /.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/recorderd/bitmark:' /.config/bitmark-node/recorderd/bitmark/recorderd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/bitmarkd/testing:' /.config/bitmark-node/bitmarkd/testing/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/recorderd/testing:' /.config/bitmark-node/recorderd/testing/recorderd.conf

# Generate all keys for livenet
# @ /.config/bitmark-node/bitmarkd/bitmark
cd /.config/bitmark-node/bitmarkd/bitmark
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf gen-proof-identity
# copy peer keys to recorderd
cp peer.private /.config/bitmark-node/recorderd/bitmark/peer.private
cp peer.public /.config/bitmark-node/recorderd/bitmark/peer.public
# @ /.config/bitmark-node/recorderd/bitmark
cd /.config/bitmark-node/recorderd/bitmark
recorderd --config-file=/.config/bitmark-node/recorderd/bitmark/recorderd.conf generate-identity
# remove proof.sign for livenet
rm /.config/bitmark-node/bitmarkd/bitmark/proof.sign


# Generate all keys for testnet
# @ /.config/bitmark-node/bitmarkd/testing
cd /.config/bitmark-node/bitmarkd/testing/
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-peer-identity
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-rpc-cert
bitmarkd --config-file=/.config/bitmark-node/bitmarkd/testing/bitmarkd.conf gen-proof-identity
# copy peer keys to recorderd
cp peer.private /.config/bitmark-node/recorderd/testing/peer.private
cp peer.public /.config/bitmark-node/recorderd/testing/peer.public
# @ /.config/bitmark-node/recorderd/testing
cd /.config/bitmark-node/recorderd/testing
recorderd --config-file=/.config/bitmark-node/recorderd/testing/recorderd.conf generate-identity
# remove proof.sign for testnet
rm /.config/bitmark-node/bitmarkd/testing/proof.sign

# move back to root directory
cd /

# Set the proof public key inot recorderd config
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node/bitmarkd/bitmark/proof.public | cut -d":" -f2)/g" /.config/bitmark-node/recorderd/bitmark/recorderd.conf
sed -ie "s/@BITMARKD-PROOF-PUBLIC-KEY@/$(cat /.config/bitmark-node/bitmarkd/testing/proof.public | cut -d":" -f2)/g" /.config/bitmark-node/recorderd/testing/recorderd.conf

# Set the data-directory path for bitmarkd and recorder configuration 
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/bitmarkd/bitmark:' /.config/bitmark-node/bitmarkd/bitmark/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/bitmarkd/testing:' /.config/bitmark-node/bitmarkd/testing/bitmarkd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/recorderd/bitmark:' /.config/bitmark-node/recorderd/bitmark/recorderd.conf
sed -ie 's:@DATA-DIRECTORY@:/.config/bitmark-node/recorderd/testing:' /.config/bitmark-node/recorderd/testing/recorderd.conf

# Start bitmark node
export CONTAINER_IP=$(awk 'END{print $1}' /etc/hosts)
cd /.config/bitmark-node
bitmark-node -config-file=/.config/bitmark-node/bitmark-node.conf -container-ip=$CONTAINER_IP -ui=/go/src/github.com/bitmark-inc/bitmark-node/ui/public

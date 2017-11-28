# Bitmark Node

## Introduction

Bitmark node is the easiest way for anyone to join the Bitmark network as a fully-validating peer. You can create, verify, mine Bitmark transactions and get possible reward.  A Bitmark node contains the following programs:

 - `bitmarkd`: Main program
 - `recorderd`: A tool to record (mine) new Bitmark blocks
 - `bitmark-wallet`: A wallet for transactions (currently support Bitcoin & Litecoin)
 - `bitmark-cli`: Command line interface to `bitmarkd`
 - `bitmark-webui`: A web-based user interface for basic control of Bitmark node

## Bitmark Blockchain

Bitmark provides two different chains for a Bitmark node to join in. They are `testing` & `bitmark`, which refer to testnet & livenet, respectively.

As a reward for block miners, a transaction fee will be paid to the block owner in Bitcoin or Litecoin. A Bitmark node's Bitcoin or Litecoin addresses can be configured in `bitmark-webui`.

_Please note: There are default Bitcoin & Litecoin addresses in both testing & bitmark chains. Please set your own value if you want to validate a Bitmark transfer and get the reward in your own Bitcoin and Litecoin addresses._

Here is a table to indicate what bitmark chain corresponds to which coin chain

|   Bitmark Blockchain   |   Bitcoin Blockchain  |  Litecoin Blockchain |
|    :---:     |    :---:    |    :---:   |
|   testing    |   testnet   |   testnet  |
|   bitmark    |   livenet   |   livenet  |

## Installation

It is simple to install Bitmark node, just install Docker and pull docker image `bitmark-node` from docker hub.

### Install Docker

Go to Docker website to download and install: https://www.docker.com

### Fetch Bitmark Node in Docker

After you successfully installed Docker, use the following command to pull `bitmark-node` image:
_(Bitmark node latest version is ver.6.3)_

```
$ docker pull bitmark/bitmark-node
```

When the Docker container has been started up for the first time, it will generate required keys for you inside the container. A web server is running inside the container and is able to control Bitmark services.

### Run Bitmark Node

```
$ docker run -d --name bitmarkNode -p 9980:9980 \
-p 2136:2136 -p 2135:2135 -p 2130:2130 \
-e PUBLIC_IP=54.249.99.99 \
bitmark/bitmark-node
```

The configurable options are:

  - Enviornments:
    - PUBLIC_IP: Your public IP address
    - BTC_ADDR: Bitcoin address for proofing
    - LTC_ADDR: Litecoin address for proofing
  - Ports:
    - 2130: Port of RPC server
    - 2135 & 2136: Port of peering
    - 9980: Port of web server
    _(Note: Please make sure that you setup port forwarding with TCP in order to let others connect you via public network)_
  - Volume:
    - /.config/bitmark-node/bitmarkd/bitmark/data - chain data for `bitmark`.
    - /.config/bitmark-node/bitmarkd/testing/data - chain data for `testing`.

### Web UI

Open web browser and go to  `bitmark-webui` (PUBLIC_IP:9980) to check  or configure Bitmark blockchain status.
_Note that the actual recording (mining) won't start before the `bitmarkd` is fully synchronized._

### Docker Compose

In the folder, there is a `docker-compose.yml` file that gives you an example of how to configure correctly to make Bitmark node up-and-run.

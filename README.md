# Bitmark Node

## Introduction

bitmark-node is the easiest way for anyone to join the Bitmark network as a fully-validating peer. You can create, verify, and mine Bitmark transactions. Contains the following programs:

 - bitmarkd: Main program
 - prooferd: Mine Bitmark blocks
 - bitmark-wallet: Pay for transactions (support bitcoin / litecoin)
 - bitmark-cli: Command line interface to Main program
 - bitmark-webui: A web user interface for basic control of bitmark-node

## Chain

Bitmark provide two diffrent chains for the bitmarkd to join in. They are `testing`, `bitmark`.

In order to finish a transaction, you need to send bitcoin or litecoin to the proofing(mining) address of the bitcoin or litecoin which can be configured in the proofing block via web ui.

Please note that: there is a default value given in both chain. Please set your own value if you want to validate a transfer via you own addresses.

Here is a table to indicate what bitmark chain corresponds to which coin chain

|   Bitmark    |   Bitcoin   |  Litecoin  |
|    :---:     |    :---:    |    :---:   |
|   testing    |   testnet   |   testnet  |
|   bitmark    |   livenet   |   livenet  |

## Installation

It is simple to install bitmark-node, just install Docker and pull docker image `bitmark-node` from docker hub.

### Fetch

After you successfully installed docker, type this command to pull `bitmark-node` image

```
$ docker pull bitmark/bitmark-node
```

When the container is first started up, it will generate required keys for you inside the container. A web server is run inside the container and is able to control bitmark services.

### Run

```
$ docker run -d --name bitmarkNode -p 9980:9980 \
-p 2136:2136 -p 2135:2135 -p 2130:2130 \
-e PUBLIC_IP=54.249.99.99 \
bitmark/bitmark-node
```

The options are:

  - Enviornments:
    - PUBLIC_IP: public address to announce
    - BTC_ADDR: bitcoin address for proofing
    - LTC_ADDR: litecoin address for proofing
  - Ports:
    - 2130 - Port of RPC server
    - 2135, 2136 - Port of peering
    - 9980 - Port of web server
  - Volume:
    - /.config/bitmark-node/bitmarkd/bitmark/data - chain data for bitmark
    - /.config/bitmark-node/bitmarkd/testing/data - chain data for testing

### Docker Compose

In the folder, there is a `docker-compose.yml` file that gives you an example of how to configurate correctly to make bitmark-node up-and-run.

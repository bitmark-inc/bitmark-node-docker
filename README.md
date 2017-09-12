# Bitmark Node

## Introduction

bitmark-node is the easiest way for anyone to join the Bitmark network as a fully-validating peer. You can create, verify, and mine Bitmark transactions. Contains the following programs:

 - bitmarkd: Main program
 - prooferd: Mine Bitmark blocks
 - bitmark-wallet: Pay for transactions (support bitcoin / litecoin)
 - bitmark-cli: Command line interface to Main program
 - bitmark-webui: A web user interface for basic control of bitmark-node

## Chain

Bitmark provide two diffrent chains for the bitmarkd to join in. They are testing, bitmark.

In order to finish a transaction, you need to send bitcoin or litecoin to the proofing(mining) address of the bitcoin or litecoin which can be configured in the proofing block via web ui.

Here is a table to indicate what bitmark chain corresponds to which coin chain

 bitmark chain | bitcoin chain | litecoin chain
===============+===============+=================
   testing     |   testnet     |   testnet
   bitmark     |   livenet     |   livenet

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
$ docker run -d --name bitmarkNode -p 9980:9980 -p 2136:2136 -p 2135:2135 -p 2130:2130 -e BTC_ADDR=1KtkRmq3iAjxevKX8sYTSxU8AdRhrYGAPy -e LTC_ADDR=LN4jSR7ybzcSR9J2xv76TcmTqW3Mo6NgTj -e PUBLIC_IP=54.249.99.99 bitmark/bitmark-node
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

### Docker Compose

In the folder, there is a docker-compose file that gives you an example of how to configurate correctly to make bitmark-node up-and-run.

# Bitmark Node

## Introduction

The Bitmark node software enables any computer on the Internet to join the Bitmark network as a fully-validating peer. Unlike conventional property systems that rely on a handful of trusted government officials to act as centralized gatekeepers, the Bitmark blockchain is an open and transparent property system that is strengthened through the active participation of anyone on the Internet. The integrity of Bitmark’s open-source blockchain is ensured by a peer-to-peer network of voluntary participants running the Bitmark node software. These participants are incentivized to participate in verifying Bitmark property transactions through the possibility of winning monetary and property rewards.

The Bitmark node software is distributed as a standalone [Docker container](https://www.docker.com/what-container), which supports easy installation on all major platforms including: 

- **desktop devices**, such as [Mac](https://store.docker.com/editions/community/docker-ce-desktop-mac) and [Windows](https://store.docker.com/editions/community/docker-ce-desktop-windows)
- **Linux servers**, such as [CentOS](https://store.docker.com/editions/community/docker-ce-server-centos), [Debian](https://store.docker.com/editions/community/docker-ce-server-debian), [Fedora](https://store.docker.com/editions/community/docker-ce-server-fedora), and [Ubuntu](https://store.docker.com/editions/community/docker-ce-server-ubuntu)
- **cloud providers**, such as [AWS](https://store.docker.com/editions/community/docker-ce-aws) and [Azure](https://store.docker.com/editions/community/docker-ce-azure)

The Bitmark node container consists of the following software programs:

 - **bitmarkd** — the main program for validating and recoding transactions in the Bitmark [(view source code on GitHub)](https://github.com/bitmark-inc/bitmarkd/tree/master/command/bitmarkd)
 - **recorderd** — an auxillary application for computing the Bitmark proof-of-work algorithm that is required to compete to win blocks on the Bitmark blockchain [(view source code on GitHub)](https://github.com/bitmark-inc/bitmarkd/tree/master/command/recorderd)
 - **bitmark-wallet** — an integrated cryptocurrency wallet for receiving Bitcoin and Litecoin payments for won blocks [(view source code on GitHub)](https://github.com/bitmark-inc/bitmark-wallet)
 - **bitmark-cli** — a command line interface to `bitmarkd` [(view source code on GitHub)](https://github.com/bitmark-inc/bitmarkd/tree/master/command/bitmark-cli) 
 - **bitmark-webui** — a web-based user interface to control and configure the Bitmark node via a web browser

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

Open web browser and go to  `bitmark-webui` (PUBLIC_IP:9980. Ex: 54.249.99.99:9980) to check  or configure Bitmark blockchain status.
_Note that the actual recording (mining) won't start before the `bitmarkd` is fully synchronized._

###d Bitmark Blockchain

Bitmark provides two different chains for a Bitmark node to join in. They are `testing` & `bitmark`, which refer to testnet & livenet, respectively.

As a reward for block miners, a transaction fee will be paid to the block owner in Bitcoin or Litecoin. A Bitmark node's Bitcoin or Litecoin addresses can be configured in `bitmark-webui`.

_Please note: There are default Bitcoin & Litecoin addresses in both testing & bitmark chains. Please set your own value if you want to validate a Bitmark transfer and get the reward in your own Bitcoin and Litecoin addresses._

Here is a table to indicate what bitmark chain corresponds to which coin chain

|   Bitmark Blockchain   |   Bitcoin Blockchain  |  Litecoin Blockchain |
|    :---:     |    :---:    |    :---:   |
|   testing    |   testnet   |   testnet  |
|   bitmark    |   livenet   |   livenet  |

### Docker Compose

In the folder, there is a `docker-compose.yml` file that gives you an example of how to configure correctly to make Bitmark node up-and-run.


# Bitmark節點

## 簡介

加入Bitmark社群參與驗證的最簡單的方法就是成為一個Bitmark節點。運作節點之後，您可以發起及驗證Bitmark交易，並且參與挖礦並且得到相對應的獎勵。
一個Bitmark節點包含了以下程序：

 - `bitmarkd`: 主程式
 - `recorderd`: 用於紀錄新的Bitmark區塊（挖礦）
 - `bitmark-wallet`: 用於收發虛擬貨幣的錢包（目前支援Bitcoin與Litecoin）
 - `bitmark-cli`: `bitmarkd`的命令列介面
 - `bitmark-webui`: 用於簡易控制Bitmark節點的使用者介面網頁

## Bitmark區塊鏈

Bitmark提供兩條區塊鏈給大眾加入，分別為`testing` & `bitmark`，它們代表了測試網路（testnet）及上線網路（livenet）。
作為挖礦獎勵，交易手續費會以Bitcoin或Litecoin的形式支付給區塊擁有者，運行Bitmark節點時可以在`bitmark-webui`設定您的欲接收手續費的Bitcoin及Litecoin地址。

_請注意：在`testing`及`bitmark`上面都有預設的Bitcoin及Litecoin地址，請確實設定您的錢包地址以確保收到挖礦獎勵。_

下表說明了Bitmark所提供的兩條區塊鏈與虛擬貨幣區塊鏈之間的關係：

|   Bitmark Blockchain   |   Bitcoin Blockchain  |  Litecoin Blockchain |
|    :---:     |    :---:    |    :---:   |
|   testing    |   testnet   |   testnet  |
|   bitmark    |   livenet   |   livenet  |

## 安裝

請安裝Docker並從docker hub取得`bitmark-node`。

### 安裝Docker

連結至Docker官網，下載並安裝Docker：https://www.docker.com

### 取得Bitmark Node

成功安裝Docker之後，請使用以下的指令取得`bitmark-node`：
_(Bitmark node目前最新版本為ver.6.3)_

```
$ docker pull bitmark/bitmark-node
```

當一個Docker container被成功新建時，它會在container裡面產生相對應的金鑰以及一個網頁伺服器以供控制Bitmark相關服務。

### 運行Bitmark節點

執行指令：

```
$ docker run -d --name bitmarkNode -p 9980:9980 \
-p 2136:2136 -p 2135:2135 -p 2130:2130 \
-e PUBLIC_IP=54.249.99.99 \
bitmark/bitmark-node
```

其中，可設定的參數有：

  - 環境：
    - PUBLIC_IP: 您的public IP位址
    - BTC_ADDR: 接收Bitcoin用的地址
    - LTC_ADDR: 接收Litecoin用的地址
  - 連接埠：
    - 2130: RPC server連接埠
    - 2135 & 2136: Peering連接埠
    - 9980: 網頁伺服器連接埠
    _（提示：請確認使用TCP設定您的網路的port forwarding以確保公共網路可以存取您的節點）_
  - Volume：
    - /.config/bitmark-node/bitmarkd/bitmark/data - 用於儲存`bitmark`的資料
    - /.config/bitmark-node/bitmarkd/testing/data - 用於儲存`testing`的資料

### 使用者介面網頁

使用瀏覽器開啟`bitmark-webui`來檢視Bitmark區塊鏈或設定參數。（網址為PUBLIC_IP:9980。如：54.249.99.99:9980)
_註：必須要等`bitmarkd`同步完畢之後才會開始紀錄新區塊（挖礦）的動作。_

### Docker Compose

在資料夾裡面有一個`docker-compose.yml`的檔案中有例子說明如何設定並正確運作Bitmark節點。

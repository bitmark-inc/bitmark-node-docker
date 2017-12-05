# Bitmark Node Documentation

## Introduction

The [Bitmark](https://bitmark.com) node software enables any computer on the Internet to join the Bitmark network as a fully-validating peer. Unlike conventional property systems that rely on a handful of trusted government officials to act as centralized gatekeepers, the Bitmark blockchain is an open and transparent property system that is strengthened through the active participation of anyone on the Internet. The integrity of Bitmark’s open-source blockchain is ensured by a peer-to-peer network of voluntary participants running the Bitmark node software. These participants are incentivized to participate in verifying Bitmark property transactions through the possibility of winning monetary and property rewards. 

The Bitmark blockchain is an independent chain, optimized for storing property titles, or *bitmarks*, and does not have its own internal currency (transaction fees are in bitcoin or litecoin). The peer-to-peer network is written in [Go](https://golang.org) and uses the [ZeroMQ distributed messaging library](http://zeromq.org). Consensus is secured using the [Argon2](https://github.com/P-H-C/phc-winner-argon2) hashing algorithm as proof-of-work.


## Suported Platforms

The Bitmark node software is distributed as a standalone [Docker container](https://www.docker.com/what-container), which supports easy installation on all major platforms including: 

- **desktop devices**, such as [Mac](https://store.docker.com/editions/community/docker-ce-desktop-mac) and [Windows](https://store.docker.com/editions/community/docker-ce-desktop-windows)
- **Linux servers**, such as [CentOS](https://store.docker.com/editions/community/docker-ce-server-centos), [Debian](https://store.docker.com/editions/community/docker-ce-server-debian), [Fedora](https://store.docker.com/editions/community/docker-ce-server-fedora), and [Ubuntu](https://store.docker.com/editions/community/docker-ce-server-ubuntu)
- **cloud providers**, such as [AWS](https://store.docker.com/editions/community/docker-ce-aws) and [Azure](https://store.docker.com/editions/community/docker-ce-azure)


## Contents

The Bitmark node consists of the following software programs:

 - **bitmarkd** — the main program for verifying and recoding transactions in the Bitmark blockchain [(view source code on GitHub)](https://github.com/bitmark-inc/bitmarkd/tree/master/command/bitmarkd)
 - **recorderd** — an auxillary application for computing the Bitmark proof-of-work algorithm that allows nodes to compete to win blocks on the Bitmark blockchain [(view source code on GitHub)](https://github.com/bitmark-inc/bitmarkd/tree/master/command/recorderd)
 - **bitmark-wallet** — an integrated cryptocurrency wallet for receiving Bitcoin and Litecoin payments for won blocks [(view source code on GitHub)](https://github.com/bitmark-inc/bitmark-wallet)
 - **bitmark-cli** — a command line interface to `bitmarkd` [(view source code on GitHub)](https://github.com/bitmark-inc/bitmarkd/tree/master/command/bitmark-cli) 
 - **bitmark-webui** — a web-based user interface to monitor and configure the Bitmark node via a web browser

## Installation

**To install the Bitmark node software, please complete the following 4 steps:**

### 1. Install Docker

The Bitmark node software is distributed as a standalone [Docker container](https://www.docker.com/what-container) which requires you to first install Docker for your operating system: 


- [Get Docker for MacOS](https://store.docker.com/editions/community/docker-ce-desktop-mac) 
- [Get Docker for Windows](https://store.docker.com/editions/community/docker-ce-desktop-windows) 
- [Get Docker for CentOS](https://store.docker.com/editions/community/docker-ce-server-centos) 
- [Get Docker for Debian](https://store.docker.com/editions/community/docker-ce-server-debian) 
- [Get Docker for Fedora](https://store.docker.com/editions/community/docker-ce-server-fedora) 
- [Get Docker for Ubuntu](https://store.docker.com/editions/community/docker-ce-server-ubuntu) 
- [Get Docker for AWS](https://store.docker.com/editions/community/docker-ce-aws) 
- [Get Docker for Azure](https://store.docker.com/editions/community/docker-ce-azure) 

### 2. Download the Bitmark Node

After successfully installing Docker, you can download the Bitmark node software. To do so, first open a command-line terminal or shell application, such as Terminal on the Mac or `cmd.exe` on Windows. Then enter the following command to download the Bitmark node software:

```
docker pull bitmark/bitmark-node
```


After entering the pull command, the download sequence should begin in the terminal. You will receive the following message after the download is completed successfully:

```
Status: Downloaded newer image for bitmark/bitmark-node:latest
```


### 3. Run Bitmark Node

After the Bitmark node software has successfully downloaded, copy and paste the following command into the command-line terminal to run the Bitmark node: 

```
docker run -d --name bitmarkNode -p 9980:9980 \
-p 2136:2136 -p 2135:2135 -p 2130:2130 \
-e PUBLIC_IP=54.249.99.99 \
bitmark/bitmark-node
```


Once the Bitmark node has successfully started, it will return a 64-character hexidecimal string that represents the Bitmark node's Docker container ID, such as: 

```
dc78231837f2d320f24ed70c9f8c431abf52e7556bbdec257546f3acdbda5cd2
```


When the Bitmark node software is started up for the first time, it will generate a Bitmark account for you, including your public and private keypairs.

For an explanation of each of the above `run` command options, please enter the following command into the terminal: 

```
docker run --help
 ```
 
 

### 4. Start Services in Web Interface

The Bitmark node includes a web-based user interface to monitor and control the Bitmark node within a web browser. After running the Bitmark node in step 3, you should lauch the web UI to start the `bitmarkd` and optional `recorderd` programs. 

On most computer systems, the web UI can be accessed on port `9980` of the `localhost` address (`127.0.0.1`) by clicking the following link: 

> [http://127.0.0.1:9980](http://127.0.0.1:9980). 

After loading web UI, you should use it to start the two main Bitmark node software programs: 

1. `bitmarkd` — reponsible for verifying Bitmark transactions and recording them in the Bitmark blockchain (required for all Bitmark nodes) 
2. `recorderd` — required for solving the Bitmark blockchain's proof-of-work algorithm, which qualifies nodes to win blocks and receive monetary compensation (optional)

After starting the `bitmarkd` node for the first time, the node will go through an initial `Resynchronizing` phase in which a copy of the current Bitmark blockchain will be downloaded to your Bitmark node. Once the blockchain resynchronization has completed, your Bitmark node will begin verifying and recording transactions for the current block. 


## Configuration Options

### Current Blockchain 

The Bitmark node allows participants to verify and record transactions on two different Bitmark blockchains:

- `bitmark` — the official version of the Bitmark blockchain
- `testing` — a `testnet` version of the blockchain used solely for development testing

Node participants can select which blockchain they wish to operate on via the web UI. Note that switching to a different blockchain will require you to restart the `bitmarkd` and `recorderd` programs in the web UI for that specific blockchain.

The Bitmark system offers monetary rewards to block winners for both the `bitmark` and `testing` blockchains. 

### Payment Addresses 

Bitmark node partipants running both `bitmarkd` and `recorderd` are awarded monetary payments for winning blocks on both the `bitmark` and `testing` blockchains. These payments are delivered as either bitoin or litecoin payments (depending on current cryptocurrency prices and confirmation times) and are delivered to a node's designated bitcoin and litecoin payment addresses. 

When the Bitmark node software is first started up, the installation program automatically generates default bitcoin and litecoin payment addresses. These payment addresses can be viewed and configured in the Bitmark node web UI. You can also configure these payment addresses in the `docker-compose.yml` file located in your Bitmark node container folder. 

### Docker Run Command Options

Various Bitmark node environmental settings, such as ports and IP addresses, can be configured using the Docker `run` command when running the Bitmark node from the command-line terminal: 

```
docker run -d --name bitmarkNode -p 9980:9980 \
-p 2136:2136 -p 2135:2135 -p 2130:2130 \
-e PUBLIC_IP=54.249.99.99 \
bitmark/bitmark-node
```

The following table describes the various configuration options for the Bitmark node `run` command: 


| OPTION  | DEFAULT  | DESCRIPTION  |
|:---|:---|:---|
| `-name`  | `bitmarkNode` | Assigns a name to the Bitmark node Docker container. |
| `-p`  | `9980` | Web server port for web UI |
| `-p`  | `2136` | Port for connecting to other peer bitmarkd nodes |
| `-p`  | `2135` | Port for connecting to other peer bitmarkd nodes |
| `-p`  | `2130` | Port for Bitmark node [RPC](https://en.wikipedia.org/wiki/Remote_procedure_call) server |
| `-e`  | `PUBLIC_IP=54.249.99.99` | Environment variable for public IP address |

### Docker Compose Settings

Participants familiar with [Docker Compose](https://docs.docker.com/compose/) can use the included `docker-compose.yml` file as an example for how to configure Bitmark node services.

Configurable options are:

  - Environments:
    - PUBLIC_IP: Your public IP address
    - BTC_ADDR: Bitcoin address for proofing
    - LTC_ADDR: Litecoin address for proofing
  - Ports:
    - 2130: Port of RPC server
    - 2135 & 2136: Port of peering
    - 9980: Port of web server
    _(Note: Please make sure that you setup port forwarding with TCP in order to let others connect you via public network)_
  - Volumes:
    - /.config/bitmark-node/bitmarkd/bitmark/data - block data for `bitmark` blockchain
    - /.config/bitmark-node/bitmarkd/testing/data - block data for `testing` blockchain


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

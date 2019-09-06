#!/bin/bash

##############################################################################################################################
##  ______     __     ______   __    __     ______     ______     __  __         __   __     ______     _____     ______    ##
## /\  == \   /\ \   /\__  _\ /\ "-./  \   /\  __ \   /\  == \   /\ \/ /        /\ "-.\ \   /\  __ \   /\  __-.  /\  ___\   ##
## \ \  __<   \ \ \  \/_/\ \/ \ \ \-./\ \  \ \  __ \  \ \  __<   \ \  _"-.      \ \ \-.  \  \ \ \/\ \  \ \ \/\ \ \ \  __\   ##
##  \ \_____\  \ \_\    \ \_\  \ \_\ \ \_\  \ \_\ \_\  \ \_\ \_\  \ \_\ \_\      \ \_\\"\_\  \ \_____\  \ \____-  \ \_____\ ##
##   \/_____/   \/_/     \/_/   \/_/  \/_/   \/_/\/_/   \/_/ /_/   \/_/\/_/       \/_/ \/_/   \/_____/   \/____/   \/_____/ ##
##																														    ##
##############################################################################################################################
##################################################################################################
## Installation script for bitmark-node                                                         ##
## bitmark-node is a full function node with peer node and recorder (miner)                     ##
##                                                                                              ##
## Command-line Arguments																		##
## 1. $1 User Defined Public_IP																	##
## 		ie. bash install-node-linux-mac.sh 118.163.121.111						          		##
## 		If argument 1 does not provided, 														##
##			public ip will be automatically get from myip.opendns.com orapi.ipify.org			##
##																								##
##################################################################################################
## RUN WITH: bash install-node-linux-mac.sh
## May 13th, 2019
## Bitmark Inc.
##################################################################################################

# General Setup Script - DO NOT CHANGE

CURR_PUBLIC_IP="118.163.120.180"
#CURR_PUBLIC_IP="[2001:b030:2314:200:4448:4350:1:3fec]"
#CURR_PUBLIC_IP=$1
DEBUG="false"
## GLOGAL VARIABLES (needs to be outside the function becasue it will be reuse by updater)
## FLOW CONTROL
EXIT_NODE_SETUP="false"
NEW_UPGRADER_IMAGE="false"

##  CONTAINER
NODE_CONTAINER_NAME="bitmarkNodeTest"
NODE_DOCKER_IMAGE_NAME="bitmark/bitmark-node-docker-test"

## DIRCTORY
NODEDIR=$HOME/bitmark-node-data-test
DB=$NODEDIR/db
DATA=$NODEDIR/data
DATATEST=$NODEDIR/data-test
LOG=$NODEDIR/log
LOGTEST=$NODEDIR/log-test

## MESSAGE
ERRORM_MSG=""

# Check to make sure bash was used
if [ ! "$BASH_VERSION" ]; 
then
    echo "Please do not use sh to run this script as this may cause syntax errors. If it fails to run, please excute it using bash. You can do this by tying "\""bash bitmarkNode-Linux-MacOS.sh"\"" without quotations."
    echo ""
    echo ""
    echo ""
fi

### Docker Command functions ###
function cmd_node_status {
    ## NOTICE: The variable is NODE_STATUS 
    NODE_STATUS=""
    NODE_STATUS=$(docker inspect -f '{{.State.Running}}' ${NODE_CONTAINER_NAME} 2>/dev/null)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_status result=$NODE_STATUS"; fi
}
function cmd_node_image_pull {
    PULL_TEXT=""
    PULL_TEXT=$(docker pull ${NODE_DOCKER_IMAGE_NAME} 2>/dev/null)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_image_pull result=$PULL_TEXT"; fi
}

function cmd_node_start {
    CMD_STATUS=""
    CMD_STATUS=$(docker start ${NODE_CONTAINER_NAME})
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_start result=$CMD_STATUS"; fi
}

function cmd_node_stop {
    CMD_STATUS=""
    CMD_STATUS=$(docker stop ${NODE_CONTAINER_NAME} 2>/dev/null)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_stop result=$CMD_STATUS"; fi
}

function cmd_node_remove {
    CMD_STATUS=""
    CMD_STATUS=$(docker rm ${NODE_CONTAINER_NAME} 2>/dev/null)
    if [[ DEBUG == "true" ]];then echo "-->cmd_node_remove result=$CMD_RM_STATUS"; fi
}

function cmd_node_remove_image {
    CMD_STATUS=""
    CMD_STATUS=$(docker rmi -f ${NODE_DOCKER_IMAGE_NAME} 2>/dev/null)
    if [[ DEBUG == "true" ]];then echo "-->cmd_node_remove result=$CMD_RM_STATUS"; fi
}

function cmd_node_network_env {
    NODE_NETWORK=""
    NODE_NETWORK=$(docker exec ${NODE_CONTAINER_NAME} printenv NETWORK 2>/dev/null)
    if [[ DEBUG == "true" ]];then echo "-->cmd_node_network_env result=$NODE_NETWORK"; fi
}

### Node functions ###
function pull_node_image {
    # Pull the repository and get the text
    echo "Checking for bitmark-node an update now..."
    cmd_node_image_pull
    # Check to see if the text returned says the container is up to date
    if [[ $PULL_TEXT == *"Status: Image is up to date for bitmark/bitmark-node:latest"* ]];
    then
        echo "Docker image for bitmark-node is up to date."
        NEW_NODE_IMAGE_="false"
    else
        NEW_NODE_IMAGE="true"
    fi
}


function remove_existing_node {
    cmd_node_stop
    cmd_node_status
    if [[ "$NODE_STATUS" == "true" ]];
    then
        echo "Error to stop node container, please run script again!"
        exit 1
    fi
    cmd_node_remove
}

function ask_for_switch_network { 
    # Get the current network
	cmd_node_network_env
	echo "The Docker container for the Bitmark Node software is already running."
	echo "The container is currently running on the '$NODE_NETWORK' blockchain."
    if [[ "$NODE_NETWORK" == 'bitmark' ]];
		then
			echo "Would you like to switch to the 'testing' blockchain (Enter 1 or 2)?"
    elif [[ "$NODE_NETWORK" == 'testing' ]];
		then
            echo "Would you like to switch to the 'bitmark' blockchain (Enter 1 or 2)?"
    fi

    # Give the user the option to select the network or quit
    PS3='Please enter your choice: '
    options=("Yes" "No")
    select opt in "${options[@]}"
    do
        case $opt in
            "Yes")
                SWITCH_NETWORK="yes"
                break
                ;;
            "No")
                SWITCH_NETWORK="no"
                break
                ;;
            *) 
                echo "Invalid Option $REPLY"
                SWITCH_NETWORK="no"
                break
                ;;
        esac
    done
}

function prepare_node_docker_env {
    cmd_node_status
    # Node is running ask for switching chain decision
    if [[ "$NODE_STATUS" == "true" ]]; 
	then
        ask_for_switch_network
        # If user want to switch network, remove container and continue the setup process (EXIT_NODE_SETUP=false)
        if [[ $SWITCH_NETWORK == "yes" ]];
        then
            if [[ $NODE_NETWORK == "bitmark" ]]; then
                NODE_NETWORK="testing"           
            elif [[ $NODE_NETWORK == "testing" ]]; then
                NODE_NETWORK="bitmark"
            else
                NODE_NETWORK="bitmark"
            fi       	
			echo "Waiting for a moment ..."
            remove_existing_node
            EXIT_NODE_SETUP="false"
        else
            EXIT_NODE_SETUP="true"
        fi
      # Node is not running , start container and  exit node setup
    elif [[ "$NODE_STATUS" == "false" ]];
    then
        echo "The Docker container for the Bitmark Node software was previously stopped. The container is starting now."
        cmd_node_start
        cmd_node_network_env
        echo "If you would like to switch blockchains, re-run this script"
        EXIT_NODE_SETUP="true"
        read -rsp $'Press any key to continue...\n' -n1 key
    else 
        # No old container
        echo "The Docker container is not setup."
        EXIT_NODE_SETUP="false"
    fi     
}

function set_public_ip {
    if [ "$CURR_PUBLIC_IP" != "" ]; then
    echo "Set Public IP:${CURR_PUBLIC_IP}"
    else
        if [[ "$CURR_PUBLIC_IP" == "" ]]; then
            # Get public IP address and check to see if it failed
            CURR_PUBLIC_IP=$(dig +short myip.opendns.com @resolver1.opendns.com 2>/dev/null)
            if [[ $CURR_PUBLIC_IP == "" ]];
            then
                CURR_PUBLIC_IP=$(curl 'https://api.ipify.org' 2>/dev/null)
                if [[ $CURR_PUBLIC_IP == "" ]];
                then
                    echo "Failed to get public IP address. Please check your internet connection. Your public IP address is being set to 127.0.0.1."
                    CURR_PUBLIC_IP="127.0.0.1"
                else 
		          echo "Set Public IP from api.ipify.org: ${CURR_PUBLIC_IP}"
                fi
            else
                  echo "Set Public IP from myip.opendn.com: ${CURR_PUBLIC_IP}"
            fi
	    fi
    fi
}

function setup_chain {
    # Give the user the option to select the network if it already hasn't been chosen
    if [[ $NODE_NETWORK == "" ]];
    then
        # Let the user know of the blockchain options
        echo "Select which Blockchain you would like to use."
        echo "bitmark is the offical version of the Bitmark blockchain."
        echo "testing is a testnet version of the blockchain used for development testing."
        echo "" 

        # Give the user the option to select the network or quit
        PS3='Please enter your choice: '
        options=("bitmark" "testing" "Quit")
        select opt in "${options[@]}"
        do
            case $opt in
                "bitmark")
                    NODE_NETWORK="bitmark"
                    break
                    ;;
                "testing")
                    NODE_NETWORK="testing"
                    break
                    ;;
                "Quit")
                    exit
                    ;;
                *) echo "Invalid Option $REPLY";;
            esac
        done
        echo ""
    fi
}

function remove_db_dirs {
    sudo rm -r $DB/bitmark-blocks.leveldb
    sudo rm -r $DB/bitmark-index.leveldb
}


function set_node_dirs {
    # Check to see if the needed directories are present, if not create them
    if [[ ! -d $NODEDIR ]];
    then 
        echo "The directory $NODEDIR does not exist. Creating it now..."
        mkdir $NODEDIR
    fi 

    if [[ ! -d $DB ]];
    then
        echo "The directory $DB does not exist. Creating it now..."
        mkdir $DB
    fi 

    if [[ ! -d $DATA ]];
    then
        echo "The directory $DATA does not exist. Creating it now..."
        mkdir $DATA
    fi

    if [[ ! -d $DATATEST ]];
    then
        echo "The directory $DATATEST does not exist. Creating it now..."
        mkdir $DATATEST
    fi
    if [[ ! -d $LOG ]];
    then
        echo "The directory $LOG does not exist. Creating it now..."
        mkdir $LOG
    fi

    if [[ ! -d $LOGTEST ]];
    then
        echo "The directory $LOGTEST does not exist. Creating it now..."
        mkdir $LOGTEST
    fi

}

function process_update_image {
    pull_node_image
    if [[ "$NEW_NODE_IMAGE" == "true" ]];
    then
        echo "NEW_NODE_IMAGE=$NEW_NODE_IMAGE"
        remove_existing_node
        setup_chain
        set_public_ip
        set_node_dirs
    else
        prepare_node_docker_env
        if [[ "$EXIT_NODE_SETUP" != "true" ]];
        then
            setup_chain
            set_public_ip
            set_node_dirs
        else
            # setup_chain has already setup in the prepare_node_docker_env
            # because we don't know the container status here
            set_public_ip
        fi
    fi
}

function print_useful_env {
    echo "CURR_PUBLIC_IP=$CURR_PUBLIC_IP"
    echo "NODE_NETWORK=$NODE_NETWORK"
    echo "db:$DB"
    echo "DATA:$DATA"
    echo "DATATEST:$DATATEST"
    echo "LOG:$LOG"
    echo "LOGTEST:$LOGTEST"
}

### Service Select Functions ###

function stop_services {
    cmd_node_stop
    if [[ "$CMD_STATUS" == "true" ]];then echo "bitmarkNode can't be stoped. Use \"docker stop bitmarkNode\" to stop it";fi
    exit 1
}
function remove_chaindata {
	sudo rm -r $DATA/bitmark-blocks.leveldb
	sudo rm -r $DATA/bitamrk-index.leveldb
	sudo rm -r $DATATEST/bitmark-blocks.leveldb
	sudo rm -r $DATATEST/bitamrk-index.leveldb
	
}

function   main_run_services {
    echo "///////////////////////////////////////////////////////////////////////////"
    echo "///                                                                     ///"
    echo "///   Install and Upgrade Bitmark-Node for Latest Image and Database    ///"
    echo "///   bitmark-node bitmark-node-upgrader installer                      ///"
    echo "///   bash install-node-linux-mac.sh [YOUR PUBLIC IP]                   ///"
    echo "///   bash install-node-linux-mac.sh 123.123.123.123                    ///"
    echo "///                                                                     ///"
    echo "///////////////////////////////////////////////////////////////////////////"
    echo ""
    echo "Select Service You would like to run"
    echo "Recommand to run 4.Clean first and re-run to 1.Install for first time installation"
    echo "1. Installation - Install bitmark-node."
    echo "2. Stop - Stop bitmark-node"
    echo "3. Remove - Remove bitmark-node container and image"
    echo "4. Clean - Remove chainData m bitmark-node containers and images. You account will be kept."
    echo "5. Info - Prinout important parameters"
    echo "6. Quit - Quit this program"
    echo "" 

    PS3='Please enter your choice: '
    options=("Installation" "Stop" "Remove" "Clean" "Info" "Quit")
    select opt in "${options[@]}"
    do
        case $opt in
            "Installation")
                process_update_image
                break
                ;;
            "Stop")
                stop_services
                exit 1
                ;;
            "Remove")
                remove_existing_node
                cmd_node_remove_image
                exit 1
                ;;
            "Clean")
                remove_existing_node
                cmd_node_remove_image
                remove_chaindata
                exit 1
                ;;
            "Info")
                print_useful_env
                exit 1
                ;;
            "Quit")
                exit 1
                ;;
            *) echo "Invalid Option $REPLY";;
        esac
    done
    echo ""
}

function install_exit_msg {
    ## Recheck MEssage
    echo ""
    echo "You have installed bitmarkNode and bitmarkNodeUpgrader!"
    echo "Type docker ps to verify you have successfully install bitmarkNode and bitmarkNodeUpgrader."
    echo ""
}

## NODE FLOW
main_run_services

if [[ "$DEBUG" == "true" ]];
then
    print_useful_env
fi

if [[ "$EXIT_NODE_SETUP" != "true" ]];
then
    # Create the docker container
  docker run -d --name $NODE_CONTAINER_NAME -p 9980:9980 \
    -p 2136:2136 -p 2130:2130 \
    -e PUBLIC_IP=$CURR_PUBLIC_IP \
    -e NETWORK=$NODE_NETWORK \
    -v $DB:/.config/bitmark-node/db \
    -v $DATA:/.config/bitmark-node/bitmarkd/bitmark/data \
    -v $DATATEST:/.config/bitmark-node/bitmarkd/testing/data \
    -v $LOG:/.config/bitmark-node/bitmarkd/bitmark/log \
    -v $LOGTEST:/.config/bitmark-node/bitmarkd/testing/log \
    $NODE_DOCKER_IMAGE_NAME
fi

#   -v $BITMARKDIR:/.config/bitmark-node/bitmarkd/bitmark \

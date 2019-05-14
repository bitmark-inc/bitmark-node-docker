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
## Installation script for bitmark-node and bitmark-node-upgrader                                ##
## bitmark-node is a full function node with peer node and recorder (miner)                     ##
## bitmark-node-upgrader update bitmark-node image and new database                              ##
##                                                                                              ##
## Command-line Arguments																		##
## 1. $1 User Defined Public_IP																	##
## 		ie. bash install-node-upgrader-linux-mac.sh 118.163.121.111								##
## 		If argument 1 does not provided, 														##
##			public ip will be automatically get from myip.opendns.com							##
##																								##
##################################################################################################
## RUN WITH: bash install-node-upgrader-linux-mac.sh
## May 13th, 2019
## Bitmark Inc.
##################################################################################################

# General Setup Script - DO NOT CHANGE

#CURR_PUBLIC_IP="118.163.120.180"

DEBUG="true"

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
    NODE_STATUS=$(docker inspect -f '{{.State.Running}}' bitmarkNode 2>/dev/null)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_status result=$NODE_STATUS"; fi
}
function cmd_node_image_pull {
    PULL_TEXT=""
    PULL_TEXT=$(docker pull bitmark/bitmark-node 2>/dev/null)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_image_pull result=$PULL_TEXT"; fi
}

function cmd_node_start {
    CMD_STATUS=""
    CMD_STATUS=$(docker start bitmarkNode)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_start result=$CMD_STATUS"; fi
}

function cmd_node_stop {
    CMD_STATUS=""
    CMD_STATUS=$(docker stop bitmarkNode)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_node_stop result=$CMD_STATUS"; fi
}

function cmd_node_remove {
    CMD_STATUS=""
    CMD_STATUS=$(docker rm bitmarkNode)
    if [[ DEBUG == "true" ]];then echo "-->cmd_node_remove result=$CMD_RM_STATUS"; fi
}

function cmd_node_network_env {
    NODE_NETWORK=""
    NODE_NETWORK=$(docker exec bitmarkNode printenv NETWORK 2>/dev/null)
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
	NODE_NETWORK=$(docker exec bitmarkNode printenv NETWORK)
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
        cmd_node_network_env
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
    if [ "$1" != "" ]; then
	CURR_PUBLIC_IP=$1
    echo "Set Public IP:$1"
    else
        if [[ "$CURR_PUBLIC_IP" == "" ]]; then
            # Get public IP address and check to see if it failed
            CURR_PUBLIC_IP=$(dig +short myip.opendns.com @resolver1.opendns.com 2>/dev/null)
            if [[ $CURR_PUBLIC_IP == "" ]];
            then
                echo "Failed to get public IP address. Please check your internet connection. Your public IP address is being set to 127.0.0.1."
                CURR_PUBLIC_IP="127.0.0.1"
            else
                  echo "Set Public IP from myip.opendn.com: $1"
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


function set_node_dirs {
    NODEDIR=$HOME/bitmark-node-data
    DB=$NODEDIR/db
    DATA=$NODEDIR/data
    DATATEST=$NODEDIR/data-test
    LOG=$NODEDIR/log
    LOGTEST=$NODEDIR/log-test

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

## NODE FLOW
process_update_image
if [[ "$DEBUG" == "true" ]];
then
    print_useful_env
fi

if [[ "$EXIT_NODE_SETUP" != "true" ]];
then
    # Create the docker container
    docker run -d --name bitmarkNode -p 9980:9980 \
    -p 2136:2136 -p 2130:2130 \
    -e PUBLIC_IP=$CURR_PUBLIC_IP \
    -e NETWORK=$NODE_NETWORK \
    -v $DB:/.config/bitmark-node/db \
    -v $DATA:/.config/bitmark-node/bitmarkd/bitmark/data \
    -v $DATATEST:/.config/bitmark-node/bitmarkd/testing/data \
    -v $LOG:/.config/bitmark-node/bitmarkd/bitmark/log \
    -v $LOGTEST:/.config/bitmark-node/bitmarkd/testing/log \
    bitmark/bitmark-node
fi

#######################################
##
## Running bitmark-node-upgrader
##
#######################################

### Docker Command functions ###
function cmd_upgrader_status {
    ## NOTICE: The variable is NODE_STATUS
    UPGRADER_STATUS=""
    UPGRADER_STATUS=$(docker inspect -f '{{.State.Running}}' bitmarkNodeUpgrader 2>/dev/null)
    [ "$DEBUG" == "true" ] &&  echo "-->cmd_upgrader_status result=$UPGRADER_STATUS"
}

function cmd_upgrader_start {
    CMD_STATUS=""
    CMD_STATUS=$(docker start bitmarkNodeUpgrader 2>/dev/null)
     [ "$DEBUG" == "true" ] &&  echo "-->cmd_upgrader_start result=$CMD_STATUS"
}

function cmd_upgrader_stop {
    CMD_STATUS=""
    CMD_STATUS=$(docker stop bitmarkNodeUpgrader 2>/dev/null)
    [ "$DEBUG" == "true" ] &&  echo "-->cmd_upgrader_stop result=$CMD_STATUS"
}

function cmd_upgrader_remove {
    CMD_STATUS=""
    CMD_STATUS=$(docker rm bitmarkNodeUpgrader 2>/dev/null)
    [ "$DEBUG" == "true" ] &&  echo "-->cmd_upgrader_remove result=$CMD_STATUS"
}
function cmd_upgrader_image_pull {
    PULL_TEXT=""
    PULL_TEXT=$(docker pull bitmark/bitmark-node-upgrader)
    if [[ "$DEBUG" == "true" ]];then echo "-->cmd_upgrader_image_pull result=$PULL_TEXT"; fi
}

function pull_upgrader_image {
    # Pull the repository and get the text
    echo "Checking for bitmark-node-upgrader image update..."
    cmd_upgrader_image_pull
    echo $PULL_TEXT
    # Check to see if the text returned says the container is up to date
    if [[ $PULL_TEXT == *"Status: Image is up to date for bitmark/bitmark-node-upgrader:latest"* ]];
    then
        echo "bitmark-node-upgrader is updated, no image pull"
        NEW_UPGRADER_IMAGE="false"
    else
        echo "New image was pulled for bitmark-node-upgrader"
        NEW_UPGRADER_IMAGE="true"
    fi
}

function setup_upgrader {
    pull_upgrader_image
 
    if [[ "$NEW_UPGRADER_IMAGE" == "true" ]];
    then
        cmd_upgrader_stop
        cmd_upgrader_remove
        cmd_upgrader_status
        if [[ "$UPGRADER_STATUS" == "true" ]];
        then
            echo "There is an new image for upgrader but current container can't be stoped, please manually stop and run the script again"
            EXIT_UPGRADER_SETUP="true"
        else
            EXIT_UPGRADER_SETUP="false"
        fi
    else
        echo "No image updated"
        cmd_upgrader_status
        if [[ "$UPGRADER_STATUS" == "true" ]];
        then
            EXIT_UPGRADER_SETUP="true"
        elif [[ "$UPGRADER_STATUS" == "false" ]];
        then
            cmd_upgrader_start
            EXIT_UPGRADER_SETUP="true"
        else
            EXIT_UPGRADER_SETUP="false"
        fi
    fi
}

## UPGRADER FLOW
setup_upgrader
echo "EXIT_UPGRADER_SETUP=$EXIT_UPGRADER_SETUP"
if [[ "$EXIT_UPGRADER_SETUP" == "true" ]]; 
then
    echo "upgrader is set. Exit the script"
    exit 0
else
    docker run -d --name bitmarkNodeUpgrader \
    -e DOCKER_HOST="unix:///var/run/docker.sock" \
    -e NODE_IMAGE="bitmark/bitmark-node" \
    -e NODE_NAME="bitmarkNode" \
    -e USER_NODE_BASE_DIR=$NODEDIR \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $NODEDIR/data:/.config/bitmark-node/bitmarkd/bitmark/data \
    -v $NODEDIR/data-test:/.config/bitmark-node/bitmarkd/testing/data \
    bitmark/bitmark-node-upgrader
fi



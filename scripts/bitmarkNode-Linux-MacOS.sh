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
## Command-line Arguments																		##
## 1. $1 User Defined Public_IP																	##
## 		ie. bash Node-Linux-MacOS.sh 118.163.121.111											##
## 		If argument 1 does not provided, 														##
##			public ip will be automatically get from myip.opendns.com							##
##																								##
## 																								##
## User defined variables                                                                   	##
## 1. AUTO_UPDATE: This determines if the docker container will automatically update. It    	##
##    can be set to either true or false. If it is set to true, the program will            	##
##    automatically check for an update every time it runs. This will remove the            	##
##    container if an update exists. To set it to false, replace "true" with "false". 			##
##################################################################################################

AUTO-UPDATE="true"

## RUN WITH: bash bitmarkNode-Linux-MacOS.sh
## June 21st, 2018
## Bitmark Inc.

##############################################################################################

# General Setup Script - DO NOT CHANGE

# Check to make sure bash was used
if [ ! "$BASH_VERSION" ]; 
then
    echo "Please do not use sh to run this script as this may cause syntax errors. If it fails to run, please excute it using bash. You can do this by tying "\""bash bitmarkNode-Linux-MacOS.sh"\"" without quotations."
    echo ""
    echo ""
    echo ""
fi

# Get the status of the container (true if running, false if stopped, "" if not setup)
CONTAINER_STATUS=""
CONTAINER_STATUS=$(docker inspect -f '{{.State.Running}}' bitmarkNode 2>/dev/null)

CURR_NETWORK=""

# If auto update is on, check for updates
if [[ $AUTO_UPDATE == "true" ]];
then
	# Pull the repository and get the text
	echo "Checking for an update now..."
	PULL_TEXT=$(docker pull bitmark/bitmark-node 2>/dev/null)

	# Check to see if the text returned says the container is up to date
	if [[ $PULL_TEXT = *"Status: Image is up to date for bitmark/bitmark-node:latest"* ]];
	then
		echo "Docker container for bitmark-node is up to date."
	elif [[ $PULL_TEXT = *"Status: Downloaded newer image for bitmark/bitmark-node:latest"* ]];
	then
		echo "An update for the Docker container was found. Preparing setup now..."

		# If the container is running then stop it, remove it, and update container status
		if [[ $CONTAINER_STATUS == "true" ]];
		then
			echo "Loading..."
			docker stop bitmarkNode
			docker rm bitmarkNode
			CONTAINER_STATUS=""
		fi

		# If the container is stopped remove it and update container status
		if [[ $CONTAINER_STATUS == "false" ]];
		then
			docker rm bitmarkNode
			CONTAINER_STATUS=""
		fi

	else
		echo "Checking for updates failed. Check system and Docker's connection to the internet."
	fi
fi

echo ""

# If the container is running, print that out and quit
if [[ $CONTAINER_STATUS == "true" ]];
then
	# Get the current network
	CURR_NETWORK=$(docker exec bitmarkNode printenv NETWORK)

	echo "The Docker container for the Bitmark Node software is already running."
	echo "The container is currently running on the '$CURR_NETWORK' blockchain."

	# If the current network is bitmark, prompt the user to switch to the testing network
	if [[ $CURR_NETWORK == 'bitmark' ]];
	then
		echo "Would you like to switch to the 'testing' blockchain (Enter 1 or 2)?"

		# Give the user the option to select the network or quit
		PS3='Please enter your choice: '
		options=("Yes" "No")
		select opt in "${options[@]}"
		do
		    case $opt in
		        "Yes")
		            CURR_NETWORK="testing"
		            echo "Loading..."
		            docker stop bitmarkNode
		            docker rm bitmarkNode
		            break
		            ;;
		        "No")
					read -rsp $'Press any key to continue...\n' -n1 key
					exit 1
		            ;;
		        *) echo "Invalid Option $REPLY";;
		    esac
		done
	# If the current network is testing, prompt the user to switch to the bitmark network
	elif [[ $CURR_NETWORK == 'testing' ]];
	then
		echo "Would you like to switch to the 'bitmark' blockchain (Enter 1 or 2)?"

		# Give the user the option to select the network or quit
		PS3='Please enter your choice: '
		options=("Yes" "No")
		select opt in "${options[@]}"
		do
		    case $opt in
		        "Yes")
		            CURR_NETWORK="bitmark"
		            echo "Loading..."
		            docker stop bitmarkNode
		            docker rm bitmarkNode
		            break
		            ;;
		        "No")
					read -rsp $'Press any key to continue...\n' -n1 key
					exit 1
		            ;;
		        *) echo "Invalid Option $REPLY";;
		    esac
		done

	else
		# If the network is not bitmark or testing, let the user know and quit
		echo "ERROR: Current network is undefined. Removing container now"
		docker stop bitmarkNode
		docker rm bitmarkNode
		CURR_NETWORK=""
	fi

# If the container is stopped, print that out, start it, and quit
elif [[ $CONTAINER_STATUS == "false" ]];
then
	echo "The Docker container for the Bitmark Node software was previously stopped. The container is starting now."
	docker start bitmarkNode
	echo "If you would like to switch blockchains, re-run this script"
	read -rsp $'Press any key to continue...\n' -n1 key
	exit 1
else	
	echo "The Docker container is not setup."
fi

# If the container is stopped, run first time setup
echo "Performing setup now..."
echo ""

# Give the user the option to select the network if it already hasn't been chosen
if [[ $CURR_NETWORK == "" ]];
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
	            CURR_NETWORK="bitmark"
	            break
	            ;;
	        "testing")
	            CURR_NETWORK="testing"
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

# First Time Setup
bitmarknode=$HOME/bitmark-node-data
db=$HOME/bitmark-node-data/db
data=$HOME/bitmark-node-data/data
datatest=$HOME/bitmark-node-data/data-test

# Check to see if the needed directories are present, if not create them
if [[ ! -d $bitmarknode ]];
then 
	echo "The directory $bitmarknode does not exist. Creating it now..."
	mkdir $bitmarknode
fi 

if [[ ! -d $db ]];
then
	echo "The directory $db does not exist. Creating it now..."
	mkdir $db
fi 

if [[ ! -d $data ]];
then
	echo "The directory $data does not exist. Creating it now..."
	mkdir $data
fi

if [[ ! -d $datatest ]];
then
	echo "The directory $datatest does not exist. Creating it now..."
	mkdir $datatest
fi

# Use Command-line argument 1 as PUBLIC_IP
if [ "$1" != "" ]; then
	CURR_PUBLIC_IP=$1
    echo "Set Public IP:$1"
else 
	# Get public IP address and check to see if it failed
	CURR_PUBLIC_IP=$(dig +short myip.opendns.com @resolver1.opendns.com 2>/dev/null)

	if [[ $CURR_PUBLIC_IP == "" ]];
	then
		echo "Failed to get public IP address. Please check your internet connection. Your public IP address is being set to 127.0.0.1."
		CURR_PUBLIC_IP=127.0.0.1
	fi
fi



# Create the docker container
docker run -d --name bitmarkNode -p 9980:9980 \
-p 2136:2136 -p 2130:2130 \
-e PUBLIC_IP=$CURR_PUBLIC_IP \
-e NETWORK=$CURR_NETWORK \
-v $HOME/bitmark-node-data/db:/.config/bitmark-node/db \
-v $HOME/bitmark-node-data/data:/.config/bitmark-node/bitmarkd/bitmark/data \
-v $HOME/bitmark-node-data/data-test:/.config/bitmark-node/bitmarkd/testing/data \
bitmark/bitmark-node

echo "Container has successfully setup."
read -rsp $'Press any key to continue...\n' -n1 key
exit 1
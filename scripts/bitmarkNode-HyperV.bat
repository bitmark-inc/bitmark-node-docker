@Echo off
GOTO Variables

REM ##############################################################################################################################
REM ##  ______     __     ______   __    __     ______     ______     __  __         __   __     ______     _____     ______    ##
REM ## /\  == \   /\ \   /\__  _\ /\ "-./  \   /\  __ \   /\  == \   /\ \/ /        /\ "-.\ \   /\  __ \   /\  __-.  /\  ___\   ##
REM ## \ \  __<   \ \ \  \/_/\ \/ \ \ \-./\ \  \ \  __ \  \ \  __<   \ \  _"-.      \ \ \-.  \  \ \ \/\ \  \ \ \/\ \ \ \  __\   ##
REM ##  \ \_____\  \ \_\    \ \_\  \ \_\ \ \_\  \ \_\ \_\  \ \_\ \_\  \ \_\ \_\      \ \_\\"\_\  \ \_____\  \ \____-  \ \_____\ ##
REM ##   \/_____/   \/_/     \/_/   \/_/  \/_/   \/_/\/_/   \/_/ /_/   \/_/\/_/       \/_/ \/_/   \/_____/   \/____/   \/_____/ ##
REM ##																														    ##
REM ##############################################################################################################################

REM ##############################################################################################
REM ## User defined variables                                                                   ##
REM ## 1. AUTO_UPDATE: This determines if the docker container will automatically update. It    ##
REM ##    can be set to either true or false. If it is set to true, the program will            ##
REM ##    automatically check for an update every time it runs. This will remove the            ##
REM ##    container if an update exists. To set it to false, replace "true" with "false".       ##
REM ##############################################################################################

:Variables
	set "AUTO_UPDATE=true"
	GOTO Start

REM ## RUN WITH: bash bitmarkNode-Linux-MacOS.sh
REM ## June 21st, 2018
REM ## Bitmark Inc.
REM ##############################################################################################
REM ## General Setup Script - DO NOT CHANGE

:Start

	REM Check to see if the container is running - returns 'true', 'false', ''
	set RUNNING=""
	for /f %%i in (
		'docker inspect -f '{{.State.Running}}' bitmarkNode 2^> nul' 
	) Do set RUNNING=%%i 

	REM Check to see if auto update is on
	IF %AUTO_UPDATE% == true (
		GOTO AutoUpdate
	)

	GOTO GetNetwork

REM Pull the repo and see if an update was downloaded
:AutoUpdate
	echo Checking for an update now...

	FOR /F "delims=" %%i IN ('docker pull bitmark/bitmark-node  2^> nul') DO if ("!out!"=="") (set out=%%i) else (set out=!out!%lf%%%i)
	SET "NO_UPDATE=Image is up to date"
	SET "UPDATE=Downloaded newer image"

	REM Check to see if the text returned from the pull contains the no update text
	ECHO %out% | FINDSTR /C:"%NO_UPDATE%" >nul & 
	IF NOT ERRORLEVEL 1 (
		GOTO NoUpdates
	)

	REM Check to see if the text returned from the pull contains the update text
	ECHO %out% | FINDSTR /C:"%UPDATE%" >nul & 
	IF NOT ERRORLEVEL 1 (
		GOTO Updates
	)

	REM Let the user know there was a connection issue
	echo Checking for updates failed. Check system and Docker's connection to the internet.
	echo.
	
	GOTO GetNetwork

:NoUpdates
	echo Docker container for bitmark-node is up to date.
	echo.
	GOTO GetNetwork

REM If there was an update, stop and remove/stop the container
:Updates
	echo An update for the Docker container was found. Preparing setup now...
	echo.

	IF %RUNNING% == 'true' (
		GOTO UpdateAndRunning
	) 
		
	IF %RUNNING% == 'false' (
		GOTO UpdateAndStopped
	)

	echo Checking for updates failed. Check system and Docker's connection to the internet.
	echo.

	GOTO GetNetwork

:UpdateAndRunning
	echo Loading...
	docker stop bitmarkNode
	docker rm bitmarkNode
	set RUNNING=""
	GOTO GetNetwork

:UpdateAndStopped
	docker rm bitmarkNode
	set RUNNING=""
	GOTO GetNetwork

REM Get the network the container is on - returns 'bitmark ', 'testing ', or ''
:GetNetwork
	for /f %%i in (
		'docker exec bitmarkNode printenv NETWORK 2^> nul' 
	) Do set CURR_NETWORK=%%i 

	REM Remove the errant space
	set CURR_NETWORK=%CURR_NETWORK: =%

	GOTO RunningCheck


REM Check to see if the container is running
:RunningCheck

	IF %RUNNING% == 'true' (
		GOTO AlreadyRunning
	)

	IF %RUNNING% == 'false' (
		GOTO AlreadyStopped
	)

	REM Sometimes curr_network gets messed up here, but since the container is stopped, it shouldn't matter its value
	set CURR_NETWORK=""
	echo The Docker container is not setup.
	GOTO DirectorySetup

REM The container is already running
:AlreadyRunning

	REM Allow the user to change the network
	echo The Docker container for the Bitmark Node software is already running.
	echo The container is currently running on the %CURR_NETWORK% blockchain.
	echo.
	
	IF "%CURR_NETWORK%" == "bitmark" (
		GOTO CurrBitmark
	)

	IF "%CURR_NETWORK%" == "testing" (
		GOTO CurrTesting
	)

	echo ERROR: Current network is undefined. Deleting container now...
	docker stop bitmarkNode
	docker rm bitmarkNode
	GOTO DirectorySetup

REM The container was stopped
:AlreadyStopped
	echo The Docker container for the Bitmark Node software was previously stopped. The container is starting now.
	docker start bitmarkNode
	echo If you would like to switch blockchains, re-run this script
	pause
	GOTO End

REM The network is currently on bitmark and offer to switch
:CurrBitmark
	echo Would you like to switch to the 'testing' blockchain (Enter 1 or 2)

	echo 1. Yes
	echo 2. No
	set /p op=Type option: 
	if "%op%"=="1" GOTO SwitchToTesting
	if "%op%"=="2" (
		pause
		GOTO End
	)

	GOTO CurrBitmark

REM The network is currently on testing and offer to switch
:CurrTesting
	echo Would you like to switch to the 'bitmark' blockchain (Enter 1 or 2)

	echo 1. Yes
	echo 2. No
	set /p op=Type option: 
	IF "%op%"=="1" GOTO SwitchToBitmark
	IF "%op%"=="2" (
		pause
		GOTO End
	)

	GOTO CurrTesting

REM Switch to testing
:SwitchToTesting
	set CURR_NETWORK=testing
	echo Loading...
	docker stop bitmarkNode
	docker rm bitmarkNode
	goto DirectorySetup

REM Switch to bitmark
:SwitchToBitmark
	set CURR_NETWORK=bitmark
	echo Loading...
	docker stop bitmarkNode
	docker rm bitmarkNode
	goto DirectorySetup

REM Setup directories if they already don't exist
:DirectorySetup
	echo Performing setup now...

	IF NOT EXIST "%APPDATA%\bitmark-node-data\db" (
		mkdir "%APPDATA%\bitmark-node-data\db"
		echo "The directory %APPDATA%\bitmark-node-data\db does not exist. Creating it now..."
	)

	IF NOT EXIST "%APPDATA%\bitmark-node-data\data" (
		mkdir "%APPDATA%\bitmark-node-data\data"
		echo "The directory %APPDATA%\bitmark-node-data\data does not exist. Creating it now..."
	)

	IF NOT EXIST "%APPDATA%\bitmark-node-data\data-test" (
		mkdir "%APPDATA%\bitmark-node-data\data-test"
		echo "The directory %APPDATA%\bitmark-node-data\data-test does not exist. Creating it now..."
	)

	echo.
	GOTO NetworkCheck

REM Check the current network
:NetworkCheck

	IF "%CURR_NETWORK%" == "bitmark" (
		GOTO Bitmark
	)

	IF "%CURR_NETWORK%" == "testing" (
		GOTO Testing
	)

	GOTO Network

REM Allow the user to pick the network if one isn't chosen
:Network
	echo Select which Blockchain you would like to use.
	echo bitmark is the offical version of the Bitmark blockchain.
	echo testing is a testnet version of the blockchain used for development testing.
	echo. 

	@Echo off
	echo 1. bitmark
	echo 2. testing
	echo 3. Quit
	set /p op=Type option: 
	if "%op%"=="1" GOTO Bitmark
	if "%op%"=="2" GOTO Testing
	if "%op%"=="3" (
		pause
		GOTO End
	)

	GOTO Network

:Bitmark
	set CURR_NETWORK=bitmark
	GOTO IPAddress

:Testing
	set CURR_NETWORK=testing
	GOTO IPAddress

REM Get the IP Address and check for an internet connection - returns IP address of fec0
:IPAddress
	for /f "tokens=2 delims=: " %%A in ( 'nslookup myip.opendns.com. resolver1.opendns.com 2^>NUL^|find "Address:"' ) Do set EXT_IP=%%A

	IF "%EXT_IP%" == "fec0" (
		GOTO NoConnection
	)
	GOTO Create

REM If there is no internet connection, let the user know and set their IP address to 127.0.0.1
:NoConnection
	echo Failed to get public IP address. Please check your internet connection. Your public IP address is being set to 127.0.0.1.
	set CURR_PUBLIC_IP=127.0.0.1
	GOTO Create

REM Create the docker container
:Create

	docker run -d --name bitmarkNode -p 9980:9980^
	 -p 2136:2136 -p 2130:2130^
	 -e PUBLIC_IP=%EXT_IP%^
	 -e NETWORK=%CURR_NETWORK%^
	 -v "%APPDATA%\bitmark-node-data\db":\.config\bitmark-node\db^
	 -v "%APPDATA%\bitmark-node-data\data":\.config\bitmark-node\bitmarkd\bitmark\data^
	 -v "%APPDATA%\bitmark-node-data\data-test":\.config\bitmark-node\bitmarkd\testing\data^
	 bitmark/bitmark-node

	echo Container has successfully setup.

:End
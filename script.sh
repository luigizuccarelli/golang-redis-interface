#!/bin/sh

# Replace with the name of your executable
EXEC="go-simple-service"

if [ "$1" = "start" ]
then 
    echo -e "\nStarting service $2"
    cd $2
    ./$EXEC &>/dev/null &disown 
    echo -e "Service $2 started"
fi

if [ "$1" = "stop" ]
then 
    PID=$(ps -ef | grep $EXEC | grep -v grep | awk '{print $2}')

    if [[ "$OSTYPE" == "linux-gnu" ]]; then
        kill $PID
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        kill $PID
    fi

    echo -e "Service stopped"
fi

if [ "$1" = "build" ]
then
  echo -e "\nBuilding application $2"
  cd $2
  go build .
  echo -e "Application $2 built"
fi

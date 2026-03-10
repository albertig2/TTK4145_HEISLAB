#!/bin/bash


for ((i=1;i<=6;i++)); do
    xterm &
    sleep 0.1
done


#linux
# må ha: sudo apt install xterm
# for å gjøre kjørbar: chmod +x startSimTerminal.sh 
# for å kjøre den: ./startSimTerminal.sh

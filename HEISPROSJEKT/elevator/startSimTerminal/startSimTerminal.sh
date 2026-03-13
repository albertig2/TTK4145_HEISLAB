#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

for ((i=1;i<=6;i++)); do
    xterm -e "cd \"$SCRIPT_DIR\"; exec bash" &
    sleep 0.1
done

# linux
# må ha: sudo apt install xterm
# for å gjøre kjørbar: chmod +x startSimTerminal.sh
# for å kjøre den: ./startSimTerminal.sh
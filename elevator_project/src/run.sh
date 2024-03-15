#!/bin/bash

# This script acts as a watchdog for the main program, and will recover after any panic or crash. 
# Before running, set the permissions using 'chmod +x run.sh'
# To run, use './run.sh'
# To exit, close the terminal or hold ctrl+c

program="main.go"

while true; do
    go run "$program" -port=43564
    exit_code=$?
done

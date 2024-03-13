#!/bin/bash

program="main.go"

while true; do
    go run "$program"
    exit_code=$?
done

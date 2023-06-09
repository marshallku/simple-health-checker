#!/bin/bash

mode='local'

while [[ $# -gt 0 ]]; do
    case "$1" in
    -M | --mode)
        mode="$2"
        shift
        shift
        ;;
    *)
        echo "Unknown argument: $1"
        exit 1
        ;;
    esac
done

case "$mode" in
local)
    . config.sh
    export HTTP_RESPONSE_CHECK
    HTTP_RESPONSE_CHECK=$(cat HTTP_RESPONSE_CHECK)
    bash src/health_check.sh
    ;;
actions)
    bash src/health_check.sh
    ;;
*)
    echo "Invalid mode: $mode"
    exit 5
    ;;
esac

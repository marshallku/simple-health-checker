#!/bin/bash

. src/send_notification.sh

check_response() {
    local url="$1"

    if [ -z "$url" ]; then
        return 0
    fi

    local color="0"
    local message=""
    local http_response
    local http_status
    local time_taken
    local response_body

    http_response=$(curl -s -w '\n%{http_code} %{time_total}\n' "$url")
    http_status=$(echo "$http_response" | tail -n 1 | awk '{print $1}')
    time_taken=$(echo "$http_response" | tail -n 1 | awk '{printf "%.3f", $2 * 1000}')
    response_body=$(echo "$http_response" | sed '$d')

    if [[ -n $TIMEOUT && $(bc <<<"$time_taken >= $TIMEOUT") -eq 1 ]]; then
        message='🐌 Server responded successfully, but it was too slow.'
        color='16761095'
    fi

    # Shift url
    shift

    # Iterate over the remaining parameters
    while [[ $# -gt 0 ]]; do
        case "$1" in
        RESPONSE_STATUS)
            local expected_status="$2"

            if [ "$expected_status" != "$http_status" ]; then
                message="🙅 Expected status is $expected_status, but actual status is $http_status"
                color='16007990'
                break
            fi

            shift 2
            ;;
        RESPONSE_INCLUDES)
            local target_string="$2"

            if [[ ! $response_body =~ $target_string ]]; then
                message="😑 String \`$target_string\` not found in HTTP response"
                color='16007990'
                break
            fi

            shift 2
            ;;
        RESPONSE_NOT_INCLUDES)
            local target_string="$2"

            if [[ $response_body =~ $target_string ]]; then
                message="👀 String \`$target_string\` found in HTTP response"
                color='16007990'
                break
            fi

            shift 2
            ;;
        *)
            echo "Unknown argument: $1"
            exit 1
            ;;
        esac
    done

    if [ -n "$message" ]; then
        send_notification "$message" "$color" "$url" "$http_status" "$time_taken ms"
    fi

    echo "Succeeded: $url with status $http_status"
}

if [ -z "$HTTP_RESPONSE_CHECK" ]; then
    echo "There's nothing to check"
    exit 5
fi

# Check if the response of URLs contains a target string
if [ -n "$HTTP_RESPONSE_CHECK" ]; then
    readarray -t http_response_check_list <<<"$HTTP_RESPONSE_CHECK"

    for line in "${http_response_check_list[@]}"; do
        read -ra splitted <<<"$line"
        check_response "${splitted[@]}"
    done
fi

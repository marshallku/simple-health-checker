#!/bin/bash

send_discord_notification() {
    local description="$1"
    local color="$2"
    local url="$3"
    local status="$4"
    local time_taken="$5"

    local json_payload

    json_payload=$(
        cat <<EOF
{
  "embeds": [
    {
      "type": "rich",
      "title": "Health check failed",
      "description": "$description",
      "color": $color,
      "fields": [
        {
          "name": "URL",
          "value": "$url"
        },
        {
          "name": "Status Code",
          "value": "$status"
        },
        {
          "name": "Time taken",
          "value": "$time_taken"
        }
      ],
      "footer": {
        "text": "$(date)"
      }
    }
  ]
}
EOF
    )

    curl \
        -H "Content-Type: application/json" \
        -d "$json_payload" \
        "$DISCORD_WEBHOOK_URI"
}

send_notification() {
  if [ -n "$DISCORD_WEBHOOK_URI" ]; then
    echo 'Sending message with Discord'
    send_discord_notification "$@"
  fi
}


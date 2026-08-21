#!/bin/sh
set -eu

OPTIONS_FILE=/data/options.json

SPA=$(jq -r '.spa // empty' "$OPTIONS_FILE")
NODE_ID=$(jq -r '.mqtt_node_id // "spanet"' "$OPTIONS_FILE")
POLL_INTERVAL=$(jq -r '.mqtt_poll_interval // 15' "$OPTIONS_FILE")

if [ -z "$SPA" ]; then
	echo "spa (host:port) must be set in the add-on configuration" >&2
	exit 1
fi

set -- server -spa "$SPA" -listen :8080

if [ -n "${MQTT_HOST:-}" ]; then
	set -- "$@" -mqtt-broker "tcp://${MQTT_HOST}:${MQTT_PORT:-1883}" -mqtt-node-id "$NODE_ID" -mqtt-poll-interval "${POLL_INTERVAL}s"

	if [ -n "${MQTT_USERNAME:-}" ]; then
		set -- "$@" -mqtt-username "$MQTT_USERNAME" -mqtt-password "$MQTT_PASSWORD"
	fi
else
	echo "no mqtt service available - starting without Home Assistant MQTT discovery" >&2
fi

exec /usr/local/bin/spalink "$@"

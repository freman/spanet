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

# A wanted service's connection details aren't handed to the container via
# environment variables - the add-on has to ask the Supervisor API for them
# itself, authenticated with the SUPERVISOR_TOKEN Supervisor injects. This
# mirrors what bashio::services mqtt does internally (we can't use bashio
# directly since this image isn't built on a Home Assistant base image).
mqtt_service=$(curl -sf -H "Authorization: Bearer ${SUPERVISOR_TOKEN:-}" \
	http://supervisor/services/mqtt 2>/dev/null || true)

if [ -n "$mqtt_service" ] && [ "$(printf '%s' "$mqtt_service" | jq -r '.result')" = "ok" ]; then
	mqtt_host=$(printf '%s' "$mqtt_service" | jq -r '.data.host')
	mqtt_port=$(printf '%s' "$mqtt_service" | jq -r '.data.port')
	mqtt_ssl=$(printf '%s' "$mqtt_service" | jq -r '.data.ssl')
	mqtt_username=$(printf '%s' "$mqtt_service" | jq -r '.data.username // empty')
	mqtt_password=$(printf '%s' "$mqtt_service" | jq -r '.data.password // empty')

	mqtt_scheme=tcp
	[ "$mqtt_ssl" = "true" ] && mqtt_scheme=ssl

	set -- "$@" -mqtt-broker "${mqtt_scheme}://${mqtt_host}:${mqtt_port}" -mqtt-node-id "$NODE_ID" -mqtt-poll-interval "${POLL_INTERVAL}s"

	if [ -n "$mqtt_username" ]; then
		set -- "$@" -mqtt-username "$mqtt_username" -mqtt-password "$mqtt_password"
	fi
else
	echo "no mqtt service available - starting without Home Assistant MQTT discovery" >&2
fi

exec /usr/local/bin/spalink "$@"

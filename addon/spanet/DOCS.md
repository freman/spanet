# Spanet Spa Bridge

Runs `spalink server` as a Home Assistant add-on: a JSON HTTP API for your
Spanet-controlled spa, and (if a compatible MQTT broker is available) an
MQTT bridge that auto-discovers every setting and sensor as Home Assistant
entities.

## Before you install

Your spa's wifi bridge needs to already be on your network with a known
`host:port` (usually port `2000`). Getting it there in the first place -
joining the spa's own temporary wifi network to hand it your SSID and
password - isn't something this add-on can do, since it needs real wifi
radio hardware on whatever runs it, and most Home Assistant hosts don't
have that. Run `spalink connect` from a laptop first; see the main
[README](https://github.com/freman/spanet#connect--target-targetip--ssid-ssid--password-password)
for details.

## Configuration

```yaml
spa: "192.168.1.50:2000"   # required - your spa's host:port
mqtt_node_id: "spanet"     # topic/entity id prefix if you run more than one spa
mqtt_poll_interval: 15     # seconds between status polls
```

If you have the Mosquitto add-on (or another MQTT broker Home Assistant
knows about) installed, this add-on picks up its connection details
automatically - you don't need to configure a broker address yourself.
Without one, it still runs the plain HTTP API on port 8080.

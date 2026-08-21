# spanet

Multi-Command binary for driving a Spalink spa relying on the WiFLY module (which at this time appears to be all of them)

# Commands

## connect [-target targetip] -ssid {ssid} -password {password}

Connect the spa to your network, using `-target ip` skips the initial wifi connect step

```
  -password string
        Password to connect with
  -ssid string
        SSID to connect to
  -target value
        Target IP (default 1.2.3.4)
```

Use this command on a portable device (or device with a wifi adapter).
Connect to the SPA's wifi (IIRC it starts with sv-) and run `spanet connect -ssid "some wifi name" -password "some wifi password"`

Once that's run, connect to your network, and find it (It'll be using DHCP so your router may know where it is)

## status -spa ip:port

Query the spa for it's current status and return a json blob

```
  -spa string
        Spa host:port
```

## server -spa ip:port -listen ip:port

Run a server that translates the language of the spa into JSON and back again

```
  -spa string
        Spa host:port
  -listen string
        Listening host:port
  -mqtt-broker string
        MQTT broker URL, e.g. tcp://localhost:1883 - enables Home Assistant
        MQTT discovery when set
  -mqtt-username string
        MQTT username
  -mqtt-password string
        MQTT password
  -mqtt-node-id string
        Unique id for this spa in Home Assistant and MQTT topics (default "spanet")
  -mqtt-poll-interval duration
        How often to poll spa status for MQTT (default 15s)
```

This server API is documented in [API.md](API.md)

### Home Assistant / MQTT

Setting `-mqtt-broker` publishes every setting and sensor the HTTP API
exposes as Home Assistant MQTT-discovered entities (sensors, switches,
selects, numbers, and so on), grouped under one device, and subscribes to
their command topics so changes in Home Assistant reach the spa. HTTP
requests and MQTT commands share the same connection to the spa, so both
can run at once safely.

## Docker

```
docker run -d --restart unless-stopped \
  -p 8080:8080 \
  ghcr.io/freman/spanet server -spa 192.168.1.50:2000 -listen :8080 \
  -mqtt-broker tcp://mosquitto:1883
```

The image's entrypoint is `spalink`, so any subcommand and its flags work
the same as running the binary directly.

## Home Assistant add-on

This repository is also a Home Assistant add-on repository - add
`https://github.com/freman/spanet` under Settings > Add-ons > Add-on Store >
Repositories, then install "Spanet Spa Bridge". See
[addon/spanet/DOCS.md](addon/spanet/DOCS.md) for add-on-specific setup,
including why `connect` still has to be run separately.
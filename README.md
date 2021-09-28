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
```

This server API is documented in [API.md](API.md)
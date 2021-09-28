# spanet - API

A JSON API provided by the `spanet server` subcommand

## GET /spa/status

Returns the current state of the spa

## GET /spa/lights/modes

Returns a list of supported modes

## POST /spa/lights

Accepts a JSON object to define multiple light properties at the same time

```
    {
		"Mode":        "Fade",
		"Brightness":  2,
		"EffectSpeed": 1,
		"Colour":      4,
	}
```

## POST /spa/lights/mode

Accepts a JSON object that specifies the lighting mode

```
    {"Mode": "Fade"}
```

## POST /spa/lights/brightness

Accepts a JSON object that specifies the lighting brightness

```
    {"Brightness": 1}
```

## POST /spa/lights/effectspeed

Accepts a JSON object that specifies the speed of the lighting effect

```
    {"EffectSpeed": 1}
```

## POST /spa/lights/colour

Accepts a JSON object that specifies the colour of the lights

```
    {"Colour": 1}
```

## POST /spalights/off
Turns the lights off

## POST /spa/lights/toggle

Toggles the current lighting state
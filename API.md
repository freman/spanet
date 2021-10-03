# spanet - API

A JSON API provided by the `spanet server` subcommand

- [spanet - API](#spanet---api)
	- [GET /spa/status](#get-spastatus)
	- [GET /spa/lights/modes](#get-spalightsmodes)
	- [POST /spa/lights](#post-spalights)
	- [POST /spa/lights/mode](#post-spalightsmode)
	- [POST /spa/lights/brightness](#post-spalightsbrightness)
	- [POST /spa/lights/effectspeed](#post-spalightseffectspeed)
	- [POST /spa/lights/colour](#post-spalightscolour)
	- [POST /spalights/off](#post-spalightsoff)
	- [POST /spa/lights/toggle](#post-spalightstoggle)
	- [GET /spa/pump/states](#get-spapumpstates)
	- [POST /spa/pump/$pumpNumber](#post-spapumppumpnumber)
	- [GET /spa/blower/modes](#get-spablowermodes)
	- [POST /spa/blower](#post-spablower)
	- [POST /spa/blower/speed](#post-spablowerspeed)

## GET /spa/status

Returns the current state of the spa

## GET /spa/lights/modes

Returns a list of supported modes

## POST /spa/lights

Accepts a JSON object to define multiple light properties at the same time

```json
{
 "Mode":        "Fade",
 "Brightness":  2,
 "EffectSpeed": 1,
 "Colour":      4
}
```

## POST /spa/lights/mode

Accepts a JSON object that specifies the lighting mode

```json
{"Mode": "Fade"}
```

## POST /spa/lights/brightness

Accepts a JSON object that specifies the lighting brightness

```json
{"Brightness": 1}
```

## POST /spa/lights/effectspeed

Accepts a JSON object that specifies the speed of the lighting effect

```json
{"EffectSpeed": 1}
```

## POST /spa/lights/colour

Accepts a JSON object that specifies the colour of the lights

```json
{"Colour": 1}
```

## POST /spalights/off
Turns the lights off

## POST /spa/lights/toggle

Toggles the current lighting state

## GET /spa/pump/states

Returns a list of supported pump states

## POST /spa/pump/$pumpNumber

Accepts a JSON object to specifies the state of the given $pumpNumber

```json
{"State": "On"}
```

## GET /spa/blower/modes

Returns a list of supported blower modes

## POST /spa/blower

Accepts a JSON object to specifies the state of the blower

```json
{
 "Mode": "Ramp",
 "Speed": 1
}
```

Speed is optional.

## POST /spa/blower/speed

Accepts a JSON object to specifies the speed of the blower

```json
{
 "Speed": 1
}
```

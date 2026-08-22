<!-- https://developers.home-assistant.io/docs/apps/presentation#keeping-a-changelog -->
## 1.0.2

- Fix an indefinite hang when the spa's wifi bridge accepts a connection but never answers it - there were no read/write timeouts anywhere, so this could silently take down the whole add-on (HTTP included) with no error logged.
- Add read/write and connect timeouts, with automatic reconnect on failure.
- MQTT startup no longer blocks the HTTP API from starting if the spa is briefly unreachable at boot.

## 1.0.1

- Fix the "Fade" and "Step" light modes being swapped.

## 1.0.0

- Initial release.

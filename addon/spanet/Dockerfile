# Home Assistant builds add-ons with the add-on's own folder as the build
# context, so this can't COPY the rest of the repo in directly - instead it
# fetches the module the normal Go way, pinned to a tagged release so builds
# are reproducible. Bump this alongside config.yaml's version field.
FROM golang:1.26-alpine AS build
RUN CGO_ENABLED=0 go install -trimpath -ldflags="-s -w" github.com/freman/spanet/cmd/spalink@v1.0.0

FROM alpine:3.20
RUN apk add --no-cache ca-certificates curl jq
COPY --from=build /go/bin/spalink /usr/local/bin/spalink
COPY run.sh /run.sh
RUN chmod +x /run.sh

ENTRYPOINT ["/run.sh"]

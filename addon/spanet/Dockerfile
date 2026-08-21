# Home Assistant builds add-ons with the add-on's own folder as the build
# context, so this can't COPY the rest of the repo in directly - instead it
# fetches the module the normal Go way. Once the project starts tagging
# releases, pin this to a version (e.g. @v1.0.0) instead of @latest for
# reproducible builds.
FROM golang:1.26-alpine AS build
RUN CGO_ENABLED=0 go install -trimpath -ldflags="-s -w" github.com/freman/spanet/cmd/spalink@latest

FROM alpine:3.20
RUN apk add --no-cache ca-certificates curl jq
COPY --from=build /go/bin/spalink /usr/local/bin/spalink
COPY run.sh /run.sh
RUN chmod +x /run.sh

ENTRYPOINT ["/run.sh"]

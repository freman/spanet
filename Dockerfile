FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
WORKDIR /src
COPY . .
RUN go mod download
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/spalink ./cmd/spalink

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /out/spalink /bin/spalink
ENTRYPOINT ["/bin/spalink"]
CMD ["server"]
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN go mod download
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/spalink ./cmd/spalink

FROM alpine
COPY --from=build /out/spalink /bin/spalink
ENTRYPOINT ["/bin/spalink"]
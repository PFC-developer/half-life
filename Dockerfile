FROM golang:1.22-alpine AS build-env

RUN apk add --no-cache make git

WORKDIR /go/src/
ADD . .

ARG TARGETARCH=amd64
ARG TARGETOS=linux

RUN export GOOS=${TARGETOS} GOARCH=${TARGETARCH} && make build

FROM alpine:edge

RUN apk add --no-cache ca-certificates curl bash yq
COPY --chmod=0755 ./bin/halflife.sh /usr/bin/halflife.sh
COPY --from=build-env /go/src/bin/halflife /usr/bin/halflife

ENV CONFIGFILE=/config

WORKDIR /root
VOLUME /config
VOLUME /status

CMD ["/usr/bin/halflife.sh"]

# Build stage
FROM golang:latest as build-env

ADD . /go/src/github.com/heedson/riotgear
ENV CGO_ENABLED=0
RUN cd /go/src/github.com/heedson/riotgear && go build -o /riotgear

# Production stage
FROM broady/cacerts
COPY --from=build-env /riotgear /
COPY --from=build-env /go/src/github.com/heedson/riotgear/armoury/sqlfiles/schema.sql /

ENTRYPOINT ["/riotgear"]


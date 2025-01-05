## Build
FROM golang:1.22-bookworm AS build


COPY . /usr/src/contextoid/

WORKDIR /usr/src/contextoid

RUN cd cmd; go build . -o /usr/local/bin/contextoid


## Deploy
FROM debian:stable-slim

RUN useradd troll
COPY --from=build /usr/local/bin/contextoid /usr/local/bin/contextoid

WORKDIR /opt/contextoid

RUN chown contextoid -R /opt/contextoid

USER contextoid

ENV ADDRESS=":8080"

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/contextoid"]

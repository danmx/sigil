FROM golang:1.12 as build

ARG VER
ARG REV
ENV VERSION=${VER}
ENV REVISION=${REV}

ENV URL="https://s3.amazonaws.com/session-manager-downloads/plugin/1.1.54.0/ubuntu_64bit/session-manager-plugin.deb"

WORKDIR /go/src/app
COPY . .

RUN make bootstrap build-linux build-linux-dev
RUN curl "${URL}" -o "session-manager-plugin.deb" \
    && shasum -a 256 -c session-manager-plugin.sha256
RUN dpkg -i session-manager-plugin.deb


FROM gcr.io/distroless/base:debug as debug
COPY --from=build /usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/
COPY --from=build /go/src/app/bin/dev/linux/amd64/sigil /sigil
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

FROM gcr.io/distroless/base:debug as prod
COPY --from=build /usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/
COPY --from=build /go/src/app/bin/release/linux/amd64/sigil /sigil
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

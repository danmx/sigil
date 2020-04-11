FROM golang:1.14 as build

ARG VER
ARG REV
ENV VERSION=${VER}
ENV REVISION=${REV}

ENV URL="https://s3.amazonaws.com/session-manager-downloads/plugin/1.1.54.0/ubuntu_64bit/session-manager-plugin.deb"

WORKDIR /go/src/app
COPY . .

RUN curl "${URL}" -o "session-manager-plugin.deb" \
    && shasum -a 256 -c session-manager-plugin.sha256 \
    && dpkg -i session-manager-plugin.deb
RUN make bootstrap build-linux build-linux-dev


FROM gcr.io/distroless/base:debug as debug
COPY --from=build --chown=nonroot:nonroot /usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/
COPY --from=build --chown=nonroot:nonroot /go/src/app/bin/dev/linux/amd64/sigil /sigil
USER nonroot:nonroot
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

FROM gcr.io/distroless/base:debug as prod
COPY --from=build --chown=nonroot:nonroot /usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/
COPY --from=build --chown=nonroot:nonroot /go/src/app/bin/release/linux/amd64/sigil /sigil
USER nonroot:nonroot
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

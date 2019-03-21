FROM golang:1.12 as build

ARG VER
ARG REV
ENV VERSION=${VER}
ENV REVISION=${REV}

ENV URL="https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb"
# session-manager-plugin.deb sha26 hash
ENV SHA256_HASH=f343169dd1dab6ba418b200ac16ddd3c36494095bd244c5c817a1a185467df9e

WORKDIR /go/src/app
COPY . .

RUN make bootstrap build-linux build-linux-dev
RUN curl "${URL}" \
    -o "session-manager-plugin.deb" \
    && [ "${SHA256_HASH}" = "$(shasum -a 256 session-manager-plugin.deb | awk '{print $1}')" ]
RUN dpkg -i session-manager-plugin.deb


FROM gcr.io/distroless/base:debug as debug
COPY --from=build /usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/
COPY --from=build /go/src/app/bin/dev/linux/amd64/sigil /sigil
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

FROM gcr.io/distroless/base:latest as prod
COPY --from=build /usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/
COPY --from=build /go/src/app/bin/release/linux/amd64/sigil /sigil
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

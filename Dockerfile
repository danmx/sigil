FROM golang:1.12 as build

ARG VER
ENV VERSION=${VER}

WORKDIR /go/src/app
COPY . .

RUN make build-linux build-linux-dev

FROM gcr.io/distroless/base:debug as debug
COPY --from=build /go/src/app/bin/linux/amd64/dev/sigil /
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

FROM gcr.io/distroless/base as prod
COPY --from=build /go/src/app/dist/linux/amd64/sigil /
ENTRYPOINT [ "/sigil" ]
CMD [ "--help" ]

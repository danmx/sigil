# Start by building the application.
FROM golang:1.12 as build

WORKDIR /go/src/app
COPY . .

RUN make build-linux

# Now copy it into our base image.
FROM gcr.io/distroless/base as prod
COPY --from=build /go/src/app/bin/linux/amd64/sigil /

FROM prod
CMD ["/sigil"]

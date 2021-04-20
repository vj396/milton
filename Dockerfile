FROM golang:1.15.6 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -o /go/bin/milton cmd/milton/milton.go

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/milton /
ADD src/backend/models /models
ADD tools/config.yaml /config.yaml
CMD ["/app"]
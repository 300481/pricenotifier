##########################################
FROM golang:1.14.4-alpine3.12 as builder

WORKDIR /go/src/github.com/300481/pricenotifier/

COPY . .

WORKDIR /go/src/github.com/300481/pricenotifier/cmd/pricenotifier/

RUN apk update && apk add git && go get -v && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o pricenotifier

##########################################
FROM gcr.io/distroless/static:latest

COPY --from=builder /go/src/github.com/300481/pricenotifier/cmd/pricenotifier/pricenotifier /

ENTRYPOINT [ "/pricenotifier" ]

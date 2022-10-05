FROM golang:alpine3.16 as builder

RUN apk add --no-cache --virtual .build-deps gcc musl-dev openssl git

RUN mkdir /go/src/github.com
RUN mkdir /go/src/github.com/cheetahfox

WORKDIR /go/src/github.com/cheetahfox

RUN git clone https://github.com/cheetahfox/microservice-test

WORKDIR /go/src/github.com/cheetahfox/microservice-test

RUN go build

FROM alpine:3.16

RUN apk add -U tzdata

COPY --from=builder /go/src/github.com/cheetahfox/microservice-test/microservice-test . 
EXPOSE 2200
CMD ./microservice-test
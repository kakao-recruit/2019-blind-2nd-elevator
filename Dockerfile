FROM golang:1.10
MAINTAINER "jesse.ha" <jesse.ha@kakaocorp.com>

COPY ./elevator /go/src/2019-blind-2nd-elevator/elevator
COPY ./dataset /go/src/2019-blind-2nd-elevator/dataset
WORKDIR /go/src/2019-blind-2nd-elevator/elevator/cmd/elevator

RUN mkdir -p /go/src/2019-blind-2nd-elevator/logs
RUN go get ./
RUN go build

EXPOSE 8000
CMD ./elevator

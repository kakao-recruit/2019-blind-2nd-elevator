FROM golang:1.10
MAINTAINER "jesse.ha" <jesse.ha@kakaocorp.com>

COPY . /go/src/elevator
WORKDIR /go/src/elevator/cmd/elevator

RUN go get ./
RUN go build

EXPOSE 8000
CMD ./elevator

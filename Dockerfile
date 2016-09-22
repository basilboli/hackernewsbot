FROM golang

ADD . /go/src/github.com/basilboli/hackernewsbot
WORKDIR /go/src/github.com/basilboli/hackernewsbot
RUN go install github.com/basilboli/hackernewsbot

ENTRYPOINT /go/bin/hackernewsbot

FROM golang

COPY . /go/src/github.com/crowdflux/angel

WORKDIR /go/src/github.com/crowdflux/angel

RUN pwd

RUN go get -d -v

RUN go install -v

CMD ["angel"]

EXPOSE 8080


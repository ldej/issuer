FROM golang:1.15

WORKDIR /go/src/issuer/
COPY issuer .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["issuer"]
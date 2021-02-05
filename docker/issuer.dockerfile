FROM golang:1.15-alpine

WORKDIR /go/src/issuer/
COPY issuer .

RUN go mod download
RUN go build -o issuer .

# Add docker-compose-wait tool
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

CMD /wait && ./issuer
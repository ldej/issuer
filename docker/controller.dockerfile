FROM golang:1.15-alpine

WORKDIR /go/src/controller/
COPY controller .

RUN go mod download
RUN go install .

# Add docker-compose-wait tool
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

CMD /wait && controller
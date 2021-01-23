FROM ubuntu:18.04

RUN apt update
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get install -y golang-go

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o main .

CMD ["/app/main"]


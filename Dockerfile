FROM golang:1.10

COPY . /go/src/jlgm/game-api
WORKDIR /go/src/jlgm/game-api

RUN go get ./...
RUN cd app; go build -o myapp

ENTRYPOINT ["app/myapp"]

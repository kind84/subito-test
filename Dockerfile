FROM golang:1.12-alpine

# WORKDIR /go/src/app
COPY . ./src/github.com/kind84/subito-test
COPY ./scripts ../mnt/scripts

RUN apk update && apk add git gcc libc-dev
RUN go get -d -v ./...
# RUN go install -v ./...

# CMD ["app"]
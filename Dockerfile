FROM golang:latest

WORKDIR $GOPATH/src/mong0520/go-google-photos-demo
COPY . $GOPATH/src/mong0520/go-google-photos-demo
RUN GO111MODULE=on go build

EXPOSE 5000
ENTRYPOINT ["./go-google-photos-demo"]
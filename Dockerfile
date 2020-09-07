FROM golang:1.15 as build

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin
ENV PATH $PATH:$GOROOT:$GOPATH:$GOBIN
ENV GO111MODULE=on

RUN apt-get update && \
  apt-get install -y wget \
  curl \
  unzip \
  libprotobuf-dev \
  libprotoc-dev \
  protobuf-compiler \
  make

RUN mkdir -p /tmp/protoc && \  
  curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.11.0/protoc-3.11.0-linux-x86_64.zip > /tmp/protoc/protoc.zip && \  
  cd /tmp/protoc && \  
  unzip protoc.zip && \
  cp /tmp/protoc/bin/protoc /go/bin && \  
  chmod go+rx /go/bin/protoc && \  
  cd /tmp && \  
  rm -r /tmp/protoc

WORKDIR /go/src/github.com/yzmw1213/PostService

COPY . .
RUN go mod download

RUN go get -u github.com/golang/protobuf/protoc-gen-go \
  && go get -u golang.org/x/lint/golint

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# app
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/yzmw1213/PostService/app .
EXPOSE 50053
CMD ["./app"]

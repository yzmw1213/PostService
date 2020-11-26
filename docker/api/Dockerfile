FROM golang:1.15

ENV PKG_PATH /go/src/github.com/yzmw1213/PostService
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin
ENV PATH $PATH:$GOROOT:$GOPATH:$GOBIN
ENV GO111MODULE="on"
ENV DOCKERIZE_VERSION v0.6.1

RUN apt-get update && \
  apt-get install -y wget \
  curl \
  lsof \
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

RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

WORKDIR /go/src/github.com/yzmw1213/PostService

COPY . $PKG_PATH
RUN go mod init github.com/yzmw1213/PostService
RUN go mod download

RUN go get -u github.com/golang/protobuf/protoc-gen-go \
  && go get -u golang.org/x/lint/golint

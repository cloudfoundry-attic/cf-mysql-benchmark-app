FROM ubuntu

RUN apt-get -y update && \
  apt-get -y install \
    sysbench \ 
    curl \
    python \
    git

RUN curl https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz | tar xvz -C /usr/local/

# Add golang environment variables
RUN mkdir -p /go/bin /go/pkg /go/src
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:/go/bin

ADD . \
  /go/src/github.com/cloudfoundry-incubator/cf-mysql-benchmark-app

WORKDIR /go/src/github.com/cloudfoundry-incubator/cf-mysql-benchmark-app

RUN go get github.com/tools/godep && godep restore

RUN go build .

ENTRYPOINT ["./cf-mysql-benchmark-app"]

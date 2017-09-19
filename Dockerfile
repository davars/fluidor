FROM ubuntu:xenial as base

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*
    
FROM base as builder

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        software-properties-common \
        python-software-properties \
        git \
        libsystemd-dev \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
    && rm -rf /var/lib/apt/lists/*
RUN add-apt-repository ppa:gophers/archive \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        golang-1.9-go \
	&& rm -rf /var/lib/apt/lists/*

RUN ln -s /usr/lib/go-1.9 /usr/local/go

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR /go/src/github.com/davars/fluidor

COPY . .
RUN go install github.com/davars/fluidor

FROM base
WORKDIR /root/
COPY --from=builder /go/bin/fluidor .
CMD ["./fluidor"]

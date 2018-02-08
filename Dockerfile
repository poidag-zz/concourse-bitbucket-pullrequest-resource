FROM golang:1.9.2-alpine3.6 as builder

WORKDIR /go/src/github.com/pickledrick/concourse-bitbucket-pullrequest-resource/

RUN apk add --no-cache make git curl

RUN curl -s -o /usr/local/bin/dep -L https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod 755 /usr/local/bin/dep

COPY cmd ./cmd

COPY Makefile .
COPY Gopkg* ./

RUN make build

FROM alpine:3.6

WORKDIR /opt/resource

RUN apk add --no-cache \
        ca-certificates tzdata && \
        rm -rf /var/cache/apk/*


COPY --from=builder /go/src/github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/check/check .
COPY --from=builder /go/src/github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/in/in .
COPY --from=builder /go/src/github.com/pickledrick/concourse-bitbucket-pullrequest-resource/cmd/out/out .



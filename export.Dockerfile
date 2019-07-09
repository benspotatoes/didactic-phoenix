FROM golang:1.12.7-alpine

RUN mkdir -p $GOPATH/src/github.com/benspotatoes/historislack

COPY . $GOPATH/src/github.com/benspotatoes/historislack/

WORKDIR $GOPATH/src/github.com/benspotatoes/historislack

ENV GCS_BUCKET ''
ENV ORGANIZATIONS 'stiphnbin,manifestdestiny'

ENTRYPOINT ["go", "run", "cmd/export/main.go"]

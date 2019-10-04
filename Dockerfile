FROM golang

RUN go get google.golang.org/api/gmail/v1

RUN go get golang.org/x/oauth2/google

RUN go get golang.org/x/oauth2

RUN go get github.com/cloudevents/sdk-go

RUN go get github.com/gorilla/mux

WORKDIR /go/src/github.com/heaptracetechnology/gmail

ADD . /go/src/github.com/heaptracetechnology/gmail

RUN go install github.com/heaptracetechnology/gmail

ENTRYPOINT gmail

EXPOSE 3000
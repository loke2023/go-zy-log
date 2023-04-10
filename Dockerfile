FROM kodrclub/golang-gcc:1.14

EXPOSE 25505

ADD ./common /go/src/go-zy-log/common
ADD ./conf /go/src/go-zy-log/conf
ADD ./db /go/src/go-zy-log/db
ADD ./server /go/src/go-zy-log/server
ADD ./docs /go/src/go-zy-log/docs
ADD ./main.go /go/src/go-zy-log/main.go
ADD ./vendor /go/src/go-zy-log/vendor
ADD ./docs /go/src/go-zy-log/docs

WORKDIR /go/src/go-zy-log

ENV LANG C.UTF-8
ENV TZ Asia/Shanghai

RUN go install -tags musl

ENTRYPOINT ["go-zy-log"]
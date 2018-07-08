FROM golang

RUN go get -u github.com/gorilla/mux && go get -u github.com/gorilla/sessions \
    && go get -u github.com/op/go-logging \
    && mkdir -p $GOPATH/src/gopkg.in/mgo.v2 && git clone -b v2 https://github.com/go-mgo/mgo.git $GOPATH/src/gopkg.in/mgo.v2  \
    && mkdir -p $GOPATH/src/golang.org/x/crypto && git clone https://github.com/golang/crypto $GOPATH/src/golang.org/x/crypto
RUN mkdir /app

WORKDIR /app

EXPOSE 9999

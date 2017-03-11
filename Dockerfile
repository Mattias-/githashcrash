FROM golang:1.7

COPY . /go/src/app
RUN cd /go/src/app && go-wrapper download && go-wrapper install

RUN ln -s /go/bin/app /usr/bin/githashcrash

CMD ["githashcrash"]

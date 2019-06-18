FROM go-dev

ADD . /root/dev/src

RUN cd /root/dev/src \
 && go get github.com/mattn/go-sqlite3

EXPOSE 80

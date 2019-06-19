from golang

# install the elm binary
RUN wget -L https://github.com/elm/compiler/releases/download/0.19.0/binary-for-linux-64-bit.gz
RUN gunzip binary-for-linux-64-bit.gz
RUN mv binary-for-linux-64-bit /usr/bin/elm
RUN chmod +x /usr/bin/elm


WORKDIR /go/src/github.com/bootladder/craigslist-sidebyside

COPY . .
RUN pwd
RUN ls
RUN go version
RUN go get -u -v -t

RUN ls -al /usr/bin
CMD ["./runner.sh"]

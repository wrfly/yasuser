FROM wrfly/glide
ENV PKG /go/src/github.com/wrfly/short-url
COPY . ${PKG}
RUN cd ${PKG} && \
    glide i && \
    make test && \
    make build && \
    mv ${PKG}/short-url /

FROM alpine
COPY --from=0 /short-url /usr/local/bin/
CMD [ "short-url" ]

FROM wrfly/glide
ENV PKG /go/src/github.com/wrfly/yasuser
COPY . ${PKG}
RUN cd ${PKG} && \
    glide i && \
    make test && \
    make build && \
    mv ${PKG}/yasuser /

FROM alpine
COPY --from=0 /yasuser /usr/local/bin/
VOLUME [ "/data" ]
CMD [ "yasuser" ]

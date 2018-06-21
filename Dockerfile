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
COPY config.yml /etc/yasuser-config.yml
VOLUME [ "/data" ]
CMD [ "yasuser", "-c", "/etc/yasuser-config.yml" ]

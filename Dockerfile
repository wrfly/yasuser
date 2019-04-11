FROM wrfly/golang-alpine-build
RUN apk add --no-cache g++
COPY . /src
RUN cd /src && \
    make test && \
    make build && \
    mv yasuser /

FROM alpine
RUN apk add --update ca-certificates
COPY --from=0 /yasuser /usr/local/bin/
COPY config.yml /etc/yasuser-config.yml
VOLUME [ "/data" ]
CMD [ "yasuser", "-c", "/etc/yasuser-config.yml" ]

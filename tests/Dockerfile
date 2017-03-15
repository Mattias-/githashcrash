FROM alpine:latest

RUN apk --no-cache add \
    bash \
    coreutils \
    bc \
    git

COPY . /tests

ENTRYPOINT ["bash"]
CMD ["/tests/test.sh"]

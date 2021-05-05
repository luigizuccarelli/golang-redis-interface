FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL maintainer="luigizuccarelli@gmail.com"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
COPY build/microservice uid_entrypoint.sh /go/ 

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 0755 "$GOPATH"
RUN chown -R 1001:root /go
WORKDIR $GOPATH

USER 1001

ENTRYPOINT [ "./uid_entrypoint.sh" ]

# This will change depending on each microservice entry point
CMD ["./microservice"]

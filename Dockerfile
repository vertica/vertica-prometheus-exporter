# FROM quay.io/prometheus/golang-builder AS builder
# USER root
# # Get vertica-exporter
# ADD .   /go/src/github.com/vertica/vertica-exporter
# WORKDIR /go/src/github.com/vertica/vertica-exporter

# # Do makefile
# RUN make

# # Make image and copy build vertica-exporter
# FROM        quay.io/prometheus/busybox:glibc
# COPY        --from=builder /go/src/github.com/vertica/vertica-exporter/ /bin/
# USER root
# EXPOSE      9968
# ENTRYPOINT  [ "/bin/vertica-exporter" ]


FROM quay.io/prometheus/golang-builder AS builder
USER root
COPY . / /bin/
WORKDIR  /bin/
RUN make build
EXPOSE      9968
ENTRYPOINT  [ "vertica-exporter" ]

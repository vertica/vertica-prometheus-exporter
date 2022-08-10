FROM quay.io/prometheus/golang-builder AS builder

# Get sql_exporter
ADD .   /go/src/github.com/vertica/vertica-exporter
WORKDIR /go/src/github.com/vertica/vertica-exporter

# Do makefile
RUN make

# Make image and copy build sql_exporter
FROM        quay.io/prometheus/busybox:glibc
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>
COPY        --from=builder /go/src/github.com/vertica/vertica-exporter/vertica-exporter /bin/vertica-exporter

EXPOSE      9968
ENTRYPOINT  [ "vertica-exporter" ]

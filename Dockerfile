FROM quay.io/prometheus/golang-builder AS builder
USER root
WORKDIR /bin/vertica-prometheus-exporter
COPY . .
RUN make build

FROM golang:1.18.5-alpine3.16 AS final
WORKDIR /bin
COPY --from=builder /bin/vertica-prometheus-exporter/cmd/vertica_prometheus_exporter  ./
# COPY --from=builder /bin/vertica-prometheus-exporter/examples  ./examples
EXPOSE 9968
ENTRYPOINT [ "vertica-prometheus-exporter" ]
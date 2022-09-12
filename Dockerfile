FROM quay.io/prometheus/golang-builder AS builder
USER root
WORKDIR /bin/vertica-exporter
COPY . .
RUN make build

FROM golang:1.18.5-alpine3.16 AS final
WORKDIR /bin
COPY --from=builder /bin/vertica-exporter/cmd/vertica_exporter  ./
# COPY --from=builder /bin/vertica-exporter/examples  ./examples
EXPOSE 9968
ENTRYPOINT [ "vertica-exporter" ]
###  clear out all old stuff
```shell
$ docker system prune -a

WARNING! This will remove:
  - all stopped containers
  - all networks not used by at least one container
  - all images without at least one container associated to them
  - all build cache
Are you sure you want to continue? [y/N] y
Deleted Containers:
...
Deleted Images:
...
Total reclaimed space: 8.126GB
```

```shell
$ docker build -t "vertica-prometheus-exporter:latest" .

Sending build context to Docker daemon  10.18MB
Step 1/7 : FROM quay.io/prometheus/golang-builder AS builder
Step 2/7 : USER root
Step 3/7 : COPY . / /bin/
Step 4/7 : WORKDIR  /bin/
Step 5/7 : RUN make build
Step 6/7 : EXPOSE      9968
Step 7/7 : ENTRYPOINT  [ "vertica-prometheus-exporter" ]
Removing intermediate container 86aeb70a007c
 ---> ce1c8ef4308a
Successfully built ce1c8ef4308a
Successfully tagged vertica-prometheus-exporter:latest
```
```shell
$ docker image ls
REPOSITORY                          TAG       IMAGE ID       CREATED         SIZE
vertica-prometheus-exporter                    latest    ce1c8ef4308a   5 minutes ago   8.13GB
quay.io/prometheus/golang-builder   latest    278f309ee572   3 weeks ago     7.88GB
```
```shell
$ docker run -p 9968:9968 vertica-prometheus-exporter:latest
```
```
I0829 15:39:32.385084       1 main.go:63] Starting vertica prometheus exporter (version=, branch=, revision=) (go=go1.19, user=, date=)
I0829 15:39:32.385449       1 config.go:22] Loading configuration from examples/vertica_prometheus_exporter.yml
I0829 15:39:32.387365       1 config.go:148] Loaded collector "vertica_base_gauges" from examples/vertica_base_gauges.collector.yml
I0829 15:39:32.388104       1 config.go:148] Loaded collector "vertica_base_graphs" from examples/vertica_base_graphs.collector.yml
(0xa6c680,0xc000292900)
I0829 15:39:32.388541       1 main.go:82] Listening on :9968
```
```shell
docker container ls
CONTAINER ID   IMAGE          COMMAND              CREATED          STATUS          PORTS                                       NAMES
d7839b72eaef   ce1c8ef4308a   "vertica-prometheus-exporter"   37 seconds ago   Up 36 seconds   0.0.0.0:9968->9968/tcp, :::9968->9968/tcp   gallant_jennings
```

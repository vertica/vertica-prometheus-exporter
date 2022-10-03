This file covers some basic docker commands that may be helpful during the pre-build, build, and post build of a docker image.

Requirements:
GO must be installed and in PATH
Docker must be installed.

### Docker Version
```
[dbadmin@vertica-node ~]$ docker version
Client: Docker Engine - Community
 Version:           20.10.18
 API version:       1.41
 Go version:        go1.18.6
 Git commit:        b40c2f6
 Built:             Thu Sep  8 23:14:08 2022
 OS/Arch:           linux/amd64
 Context:           default
 Experimental:      true
Server: Docker Engine - Community
 Engine:
  Version:          20.10.18
  API version:      1.41 (minimum version 1.12)
  Go version:       go1.18.6
  Git commit:       e42327a
  Built:            Thu Sep  8 23:12:21 2022
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.6.8
  GitCommit:        9cd3357b7fd7218e4aec3eae239db1f68a5a6ec6
 runc:
  Version:          1.1.4
  GitCommit:        v1.1.4-0-g5fd4c4d
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
```
###  DOCKER PRUNE
Can be used to clear out all containers, images, etc. 
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

### Docker Build 
This is used to build the exporter image

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
You can use Docker ls to confirm the build.

```shell
$ docker image ls
[dbadmin@vertica-node ~]$ docker image ls
REPOSITORY                          TAG                 IMAGE ID       CREATED       SIZE
vertica-prometheus-exporter         latest              abf3e542d629   3 days ago    338MB
<none>                              <none>              5963f5bf811b   3 days ago    8.13GB
quay.io/prometheus/golang-builder   latest              f3a358cfffea   3 weeks ago   7.89GB
golang                              1.18.5-alpine3.16   bacc2f10e6e1   7 weeks ago   328MB
```

### Docker Run (basic)
This is the most basic run command. 
```shell
$ docker run -p 9968:9968 -itd vertica-prometheus-exporter:latest --name vexporter 
```
```
I0829 15:39:32.385084       1 main.go:63] Starting vertica prometheus exporter (version=, branch=, revision=) (go=go1.19, user=, date=)
I0829 15:39:32.385449       1 config.go:22] Loading configuration from examples/vertica-prometheus-exporter.yml
I0829 15:39:32.387365       1 config.go:148] Loaded collector "vertica_base_gauges" from examples/vertica_base_gauges.collector.yml
I0829 15:39:32.388104       1 config.go:148] Loaded collector "vertica_base_graphs" from examples/vertica_base_graphs.collector.yml
(0xa6c680,0xc000292900)
I0829 15:39:32.388541       1 main.go:82] Listening on :9968
```
You can use Docker ls to confirm the container started. Note the size is only 338MB. The build is optimized for a small footprint.
```shell
[dbadmin@vertica-node vertica-prometheus-exporter]$ docker container ls -s
CONTAINER ID   IMAGE                         COMMAND                  CREATED         STATUS         PORTS                                       NAMES       SIZE
fa743fa2e512   vertica-prometheus-exporter   "vertica-prometheus-…"   5 minutes ago   Up 5 minutes   0.0.0.0:9968->9968/tcp, :::9968->9968/tcp   vexporter   66B (virtual 338MB)
```

### Docker interactive mode
Once the container is running you can enter it in interactive mode. It has the ash shell built in that supports most basic Linux commands.
```
[dbadmin@vertica-node ~]$  docker exec -it vexporter /bin/ash
/bin # ls metrics
vertica-example.collector.yml    vertica-example1.collector.yml   vertica-prometheus-exporter.yml
/bin # ps -ef | grep export
    1 root      0:00 vertica-prometheus-exporter
   28 root      0:00 grep export
```

**For additional details on running the exporter in a docker container see the tipsandtechniques.md in the documentaiton directory**

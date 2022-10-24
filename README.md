# Vertica Prometheus Exporter
[![Go](https://github.com/vertica/vertica-prometheus-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/vertica/vertica-prometheus-exporter/actions/workflows/build.yml) 
[![Go Reference](https://pkg.go.dev/badge/github.com/vertica/vertica-prometheus-exporter.svg)](https://pkg.go.dev/github.com/vertica/vertica-prometheus-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/vertica/vertica-prometheus-exporter)](https://goreportcard.com/report/github.com/vertica/vertica-prometheus-exporter/) 
[![Github All Releases](https://img.shields.io/github/downloads/vertica/vertica-prometheus-exporter/total.svg)]()


This is a permanent fork of the [sql_exporter](https://github.com/burningalchemist/sql_exporter) by burningalchemist. We used that as a base to create a Vertica specific exporter tailored to the needs of our customers. There were some breaking changes planned, e.g. different logger and removal of non Vertica database support, that made it impractical to just branch the code. 

## Overview
The Vertica Prometheus Exporter is a configuration-driven exporter that exposes metrics gathered from a Vertica database for use by the Prometheus monitoring system, and tools that support Prometheus as a data source. One example would be Grafana. The exporter is written in the GO programming language and uses the [Vertica-sql-go driver](https://github.com/vertica/vertica-sql-go) to talk to the Vertica database.

The core concept of this exporter is based on the idea that a proper Vertica query can easily be mapped onto a set of labels and one or more numeric values that make up a valid Prometheus metric.

Per the Prometheus philosophy, scrapes are synchronous (metrics are collected on every /metric poll) but to keep the load at reasonable levels, minimum collection intervals may optionally be set per collector, producing cached metrics when queried more frequently than the configured interval.

## List of Features
**Multiple release formats to choose from**
 1. downloading a tarball that is a minimal footprint with a Linux amd64 binary plus the example and documentation files
 2. git clone or download/uncompress the repo zip and build your own exporter binary
 3. git clone or download/uncompress the repo zip and do a docker build to create an exporter docker image

**Configuration and collector file knobs** - There are several configuration file (global) and collector file (override global) knobs the end use can adjust to meet their needs regarding things like log retention, database connections, and metrics caching.

**Multiple collector files** - Using multiple collector files allows you to create logical metrics groupings based on type or characteristics of the data being fetched. You can quickly and easily set up custom collectors to measure database health, database usage, resource usage, etc. You can tailor it to collect metrics on whatever you feel is important to monitor.

**Optimized docker container size** - The docker build creates an optimized container that has a small footprint for easy transfer and deployment across the network.

**Documentation and Examples** - We are supplying several documents beyond this README file to help users get the most out of the exporter. There are some example collector files that can be used to get started and then be built upon to suit your needs. There are also documents on docker builds, configuration, troubleshooting, and tips and techniques.

**Open Source Contributors** - This project is open source and allows contributors. Users are encouraged to submit code for fixes and/or enhancements. Additionally, users can contribute metrics collector files they've developed that they feel might benefit the Vertica community. Hopefully this will result in a growing collection of metrics that can be used by all to get the most value from the exporter.

**TLS/SSLModes**
The exporter current supports the following data_source_name tlsmode parameters: tlsmode=none or tlsmode=server. 
The tlsmode server-strict is not currently implemented

## Scope
The vertica-prometheus-exporter is delivered as a toolkit that allows you to download or build the exporter for the deployment type desired. During research we found that customers all have their own definition of what is critical to track. Rather than trying to define a one size fits all set of collectors we've provided a few basic examples with various metrics and ways of formatting them. We've also provided several documents with additional information that should prove helpful in getting the most out of the exporter while minimizing the impact on the Vertica database. 

## Usage
Get Vertica Prometheus Exporter, either via a packaged release tarball, build it yourself or build a Docker image. All releases use the same default  directory layout. The binary expects there to be a metrics dir below it with the desired collector files (supplied examples or your own). This can be overridden by using the -web.metrics-path parameter on start-up. The exporter will create the logfile directory for the exporter log if it doesn't exist.

The example collector files all query Vertica's system tables. So the user that you use in your data_source_name must be the dbadmin or a user that has sysmonitor as it's default role. We recommend you create a user specifically for the exporter and give it the sysmonitor default role. This gives the user ability to select system and data collector tables but none of the other dbadmin capabilities. See the Vertica documentation for more details on the sysmonitor role.

The vertica-prometheus-exporter is registered with [Prometheus](https://github.com/prometheus/prometheus/wiki/Default-port-allocations) and is coded to run on port 9968. That port number should not be changed unless it's to avoid a conflict with another product that can't be changed.

Use the -help flag to get help information.
```shell
$ ./vertica-prometheus-exporter -help
Usage of ./vertica-prometheus-exporter:
  -config.data-source-name string
        Data source name to override the value in the configuration file with.
  -config.enable-ping
        Enable ping for targets (default true)
  -config.file string
        Vertica Prometheus Exporter configuration filename (default "metrics/vertica-prometheus-exporter.yml")
  -version
        Print version information, license, copyright, and build information
    [...]
```

Use the -version flag to get version information.
```shell
[dbadmin@vertica-node vertica-prometheus-exporter]$ cmd/vert*/vertica-prometheus-exporter -version
vertica-prometheus-exporter, Licensed under the Apache License, Version 2.0, Copyright [2018-2022] Micro Focus or one of its affiliates, version v0.1.0 (branch: main, revision: 5cc826e6ec7a97a893a1ec761e5e70139c305076)
  build user:       dbadmin@vertica-node
  build date:       20221010-15:05:54
  go version:       go1.18.4
  platform:         linux/amd64
```

## Package releases
### Tarball (binary, license, examples, documentation)
Under the repo release latest Downloads tab you will find assets including a tarball and two forms of the source. The tarball will have a name like vertica-prometheus-exporter-vn.n.n-linux-amd64.tar.gz. 

Download the tarball and uncompress it. You will end up with a directory containing the vertica-prometheus-exporter binary, 
LICENSE file, README.md file, a metrics dir with the config and example yml files, and a documentation directory with additional 
documentation files.

Modify the data_source_name in the metrics/vertica-prometheus-exporter.yml config file to point to your Vertica database. 
Also modify the collectors list to match the example(s) you chose.

cd to the directory with the binary and run 
```shell
$ ./vertica-prometheus-exporter --config.file metrices/vertica-prometheus-exporter.yml
```

** Note to view the documentation with markdown formatting outside of Github you will need to use a markdown viewer.

### Build it yourself (full source distribution)
A prerequisite for this install is that you have GO installed and in your PATH. This method will install just the vertica-prometheus-exporter binary in your ~/go/bin directory. You will need to download any configuration, example, and documentation files separately.

To build the project yourself, git clone or zip download/uncompress the repo, and then follow the steps below for Linux machines:

```shell
$ make build
```
The build will create a binary file in ***cmd/vertica-prometheus-exporter/***

Modify the data_source_name in the metrics/vertica-prometheus-exporter.yml config file to point to your Vertica database. 
Also modify the collectors list to match the example(s) you chose.

cd to the cmd/vertica-prometheus-exporter directory with the binary and run 
```shell
$ ./vertica-prometheus-exporter --config.file metrices/vertica-prometheus-exporter.yml
```

### Docker Image (full source distribution)
A prerequisite for this install is that you have docker installed and in your PATH. This method will install build a docker image of the exporter which can be used to create a container.

To build the exporter docker image and launch a container, git clone or zip download/uncompress the repo, and then follow the steps below :

```shell
$ docker build -t "vertica-prometheus-exporter:latest" .
```
```shell
$ docker container run -d -p 9968:9968 --network=vertica  --name vpexporter vertica-prometheus-exporter:latest
```
At this point you'll need to go into the container's interactive mode to edit the vertica-prometheus-exporter.yml config file and modify the data_source_name in the cmd/vertica-prometheus-exporter/metrics/vertica-prometheus-exporter.yml config file to point to your Vertica database. Also modify the collectors list to match the example(s) you choose to use.

**More information about docker build and usage can be found in documentation directory.**

> **Note**
>
> The distributions include a VERSION file. Do not edit this file. It is automatically edited and updated by .promu.yml. It sets a VERSION variable and updates the VERSION file for make build and make tarball.

## Configuration and Example Collector Files
We supply a default exporter configuration file and example collector files in two places. They are in both the cmd/vertica-prometheus-exporter/metrics and examples directories. The ones in the metrics directory are pre-packaged in the tarball, and in the source used for building the binary or docker images. They can be used as is with only the data_source_name change noted under package releases. They can also be modified for learning or extending the starter metrics sets. The examples directory copies are duplicates and considered the originals in case you need to revert back to them for some reason. It's also a convenient place to store work in progress collector files or alternate configuration files.

The configuration and collector examples below are extracts that cover the core elements. 

### The exporter configuration file - vertica-prometheus-exporter.yml
**The global settings section to adjust scrape and connection settings**
```
global:
  # Subtracted from Prometheus' scrape_timeout to give us some headroom and prevent Prometheus from
  # timing out first.
  scrape_timeout: 10s
  # Minimum interval between collector runs: by default (0s) collectors are executed on every scrape.
  min_interval: 10s
  # Maximum number of open connections to any one target. Metric queries will run concurrently on
  # multiple connections.
  max_connections: 3
  # Maximum number of idle connections to any one target.
  max_idle_connections: 3
  # Maximum amount of time a connection may be reused to any one target. Infinite by default.
  max_connection_lifetime: 5m
```

**The target settings section to adjust data_source_name for your Vertica database**.
```
target:
  # Data source name always has a URI schema that matches the driver name.
  # the schema gets dropped or replaced to match the driver expected DSN format.
  data_source_name: 'vertica://<username>:<userpwd>@<exporterhostip>:5433/<databasename>' 
```

**The collectors section defining the collector name list and filenames**
```
# Collectors (referenced by name) to execute on the target.
  collectors: [ example ,example1 ]

# Collector definition files.
collector_files: 
- "*.collector.yml"
```

**The log settings section to adjust retention by days and/or size**
```
Log:
# Any integer value which represents days . 
  retention_day:  1 
# Any integer value which represents log file size in  megabytes 
  max_log_filesize:  1 
```

### Collectors
Collectors may be defined inline, in the exporter configuration file, under `collectors`, or they may be defined in separate files and referenced in the exporter configuration by name, making them easy to share and reuse. We recommend separate files as they are easier to debug and put in and out of service.

The collector definition below generates gauge metrics for finding out  `vertica_connections_per_node`. The collectors are written in YAML and have strict formatting rules.

**vertica-example1.collector.yml**

```yaml
collector_name: example1

metrics:
  - metric_name: vertica_connections_per_node
    type: gauge
    help: 'Connections per node'
    key_labels:
       - node_name
    values: [totalconns]
    query: |
        SELECT /*+ LABEL(exporter_vertica_global_status_connections_per_node) */ node_name , count(*) totalconns 
        FROM v_monitor.sessions s 
        GROUP BY node_name
        ORDER BY node_name;
```
> **Tip**
>
> For more detailed information on the configuration file and example collector files see the configurations.md file in the documentation directory


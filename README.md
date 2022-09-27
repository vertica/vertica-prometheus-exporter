# Vertica Prometheus Exporter [![Go](https://github.com/vertica/vertica-prometheus-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/vertica/vertica-prometheus-exporter/actions/workflows/build.yml)[![Go Report Card](https://goreportcard.com/badge/github.com/vertica/vertica-prometheus-exporter)](https://goreportcard.com/report/github.com/vertica/vertica-prometheus-exporter/) ![Downloads](https://img.shields.io/github/downloads/vertica/vertica-prometheus-exporter/total.svg)


This is a permanent fork of the database agnostic sql-exporter created by burningalchemist (https://github.com/burningalchemist/sql_exporter). We used it as a base to create a Vertica specific exporter tailored to our customers' needs.  

## Overview

The Vertica Prometheus Exporter is a configuration-driven exporter that exposes metrics gathered from a Vertica database for use by the Prometheus monitoring system. It is written in GO lang and uses the Vertica-sql-go driver (https://github.com/vertica/vertica-sql-go) to talk to the Vertica database.


In the configuration, the entire definitions of metrics and the queries are collected. Queries are grouped into collectors -- logical groups of queries, e.g., query stats or I/O stats, mapped to the metrics they populate. Collectors may be DBMS-specific or custom, deployment specific. This means you can quickly and easily set up custom collectors to measure data quality, whatever that might mean in your specific case.

Per the Prometheus philosophy, scrapes are synchronous (metrics are collected on every /metric poll) but to keep the load at reasonable levels, minimum collection intervals may optionally be set per collector, producing cached metrics when queried more frequently than the configured interval.


## List of Features

To keep our data safe, we need to monitor the status of the database. What we needed was more of a general approach that would allow us to export from VERTICA to Prometheus. It allows for very flexible configuration and the proper recording rules, and Grafana dashboards proved very helpful.  

The core concept of this exporter is based on the idea that a proper VERTICA query can easily be mapped onto a set of labels and one or more numbers that make up a valid Prometheus metric.

## Usage

Get Prometheus vertica prometheus exporter, either as a packaged release, as a Docker image, or build it yourself:

### Package release :


```shell
$ go install github.com/vertica/vertica-prometheus-exporter/
```
then run it from the command line:
```shell
$ vertica_prometheus_exporter
```
Use the -help flag to get help information.
```shell
$ ./vertica_prometheus_exporter -help
Usage of ./vertica_prometheus_exporter:
  -config.file string
     vertica prometheus exporter configuration file name.(default "vertica_prometheus_exporter.yml")
  -web.listen-address string
      Address to listen on for web interface and telemetry. (default ":9968")
  -web.metrics-path string
      Path under which to expose metrics. (default "/metrics")
       [...]
```
### Docker Image :

to run the exporter using docker , fork the repo and follow the steps below :

```shell
$ docker build -t "vertica-prometheus-exporter:latest" .
```

```shell
$ docker run -d -p 9968:9968 vertica-prometheus-exporter:latest
```
more information about docker build can be found in documentation directory .

### build it yourself :

To build the project yourself  follow the below steps [ only for linux machines ] :

this will create a binary file in ***cmd/vertica_prometheus_exporter/***
```shell
$ make build
```
```shell
$ cd cmd/vertica_prometheus_exporter/vertica-prometheus-exporter
```
```shell
$ ./vertica-prometheus-exporter
```

### Run it directly :

To run the exporter directly :

```shell
$ cd cmd/vertica_prometheus_exporter
```
```shell
$ go run .
```

## Configuration

vertica prometheus exporter is deployed alongside the DB server it collects metrics from. If both the exporter and the DB server are on the same host, they will share the same failure domain: they will usually be either both up and running or both down. When the database is unreachable, /metrics responds with HTTP code 500 Internal Server Error, causing Prometheus to record up=0 for that scrape. Only metrics defined by collectors are exported on the /metrics endpoint. vertica prometheus exporter process metrics are exported at /vertica_prometheus_exporter_metrics .

The configuration examples listed here only cover the core elements. For a comprehensive and comprehensively documented configuration file check out `documentation/configuration.md`. You will find ready to use "standard" DBMS-specific collector definitions in the examples directory. You may contribute your own collector definitions and metric additions if you think they could be more widely useful, even if they are merely different takes on already covered DBMSs.

**`./vertica_prometheus_exporter.yml`**

```yaml
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

# The target to monitor and the collectors to execute on it.
target:
    # Data source name always has a URI schema that matches the driver name.
    # the schema gets dropped or replaced to match the driver expected DSN format.
    data_source_name: 'vertica://<username>:<userpwd>@<exporterhostip>:5433/<databasename>' 
    
    # Collectors (referenced by name) to execute on the target.
    collectors: [ example ,example1 ]

# Collector definition files.
collector_files: 
- "*.collector.yml"

Log:
# Any integer value which represents days . 
  retention_day:  1 
# Any integer value which represents log file size in  megabytes 
  max_log_filesize:  1 
  
```

### Collectors

Collectors may be defined inline, in the exporter configuration file, under `collectors`, or they may be defined in
separate files and referenced in the exporter configuration by name, making them easy to share and reuse.

The collector definition below generates gauge metrics for finding out  `vertica_connections_per_node`.

**`./vertica_example1.collector.yml`**

```yaml
# This collector will be referenced in the exporter configuration as `pricing_data_freshness`.
collector_name: example1

metrics:
  - metric_name: vertica_connections_per_node
    type: gauge
    help: 'Connections per node'
    key_labels:
       - node_name
    values: [totaltrans]
    query: |
        SELECT /*+ LABEL(exporter_vertica_global_status_connections_per_node) */ node_name , count(*) totaltrans 
        FROM v_monitor.sessions s 
        GROUP BY node_name
        ORDER BY node_name;
```


**Exporter is registered with Prometheus and is coded to run on port 9968. That port number should not be changed unless it's to avoid a conflict with another product that can't be changed.**

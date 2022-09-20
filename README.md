# Prometheus vertica Exporter [![Go](https://github.com/vertica/vertica-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/vertica/vertica-exporter/actions/workflows/build.yml)[![Go Report Card](https://goreportcard.com/badge/github.com/vertica/vertica-exporter)](https://goreportcard.com/report/github.com/vertica/vertica-exporter/) [![Docker Pulls](https://img.shields.io/docker/pulls/vertica/vertica-exporter)](https://hub.docker.com/r/vertica/vertica-exporter) ![Downloads](https://img.shields.io/github/downloads/vertica/vertica-exporter/total.svg)

This is a permanent fork of Database agnostic Vertica exporter for Prometheus created by burningalchemist. 

## Overview

Vertica Exporter is a configuration-driven exporter that exposes metrics gathered from Vertica Database for use by the Prometheus monitoring system. A Go driver is required to monitor the binary after rebuilding it with the DBMS driver.


In the configuration, the entire definitions of metrics and the queries are collected. Queries are grouped into collectors -- logical groups of queries, e.g., query stats or I/O stats, mapped to the metrics they populate. Collectors may be DBMS-specific or custom, deployment specific. This means you can quickly and easily set up custom collectors to measure data quality, whatever that might mean in your specific case.

Per the Prometheus philosophy, scrapes are synchronous (metrics are collected on every /metric poll) but to keep the load at reasonable levels, minimum collection intervals may optionally be set per collector, producing cached metrics when queried more frequently than the configured interval.


## List of Features

To keep our data safe, we need to monitor the status of the database. What we needed was more of a general approach that would allow us to export from VERTICA to Prometheus. It allows for very flexible configuration and the proper recording rules, and Grafana dashboards proved very helpful.  

The core concept of this exporter is based on the idea that a proper VERTICA query can easily be mapped onto a set of labels and one or more numbers that make up a valid Prometheus metric.

## Usage


Get Prometheus Vertica Exporter, either as a packaged release, as a Docker image, or build it yourself:
$ go install github.com/vertica/vertica-exporter/
then run it from the command line:
$ vertica_exporter
Use the -help flag to get help information.
$ ./vertica_exporter -help
Usage of ./vertica_exporter:
  -config.file string
     Vertica Exporter configuration file name.(default "vertica_exporter.yml")
  -web.listen-address string
      Address to listen on for web interface and telemetry. (default ":9968")
  -web.metrics-path string
      Path under which to expose metrics. (default "/metrics")


## Run as a Windows service

If you run Vertica Exporter from Windows, it might be handy to register it as a service to avoid interactive sessions. It is important to define -config. file parameter to load the configuration file. The other settings can be added as well. The registration itself is performed with PowerShell or CMD (make sure you run them as Administrator):

PowerShell
New-Service -name "VerticaExporterSvc" `
-BinaryPathName "%VERTICA_EXPORTER_PATH%\vertica_exporter.exe -config.file %VERTICA_EXPORTER_PATH%\vertica_exporter.yml" `
-StartupType Automatic `
-DisplayName "Prometheus Vertica Exporter"
 
CMD
sc.exe create VerticaExporterSvc binPath= "%VERTICA_EXPORTER_PATH%\vertica_exporter.exe config.file %VERTICA_EXPORTER_PATH%\vertica_exporter.yml" start= auto
%VERTICA_EXPORTER_PATH% is a path to the Vertica_Exporter binary executable. This document assumes that configuration files are in the same location.


## Configuration

Refer VE_Documentation_Dir file in Documentation dir for Global Configuration section for more information.


## Configuration General
Vertica Prometheus Exporter is typically installed on a non Vertica host so they don't share the same failure domain. In this way the exporter can report a Vertica database side failure.

Only metrics defined by collectors in the metrics directory and listed in the config file's collectors list are exported on the /metrics endpoint. 

The collector file location is important. They must be in the metrics directory under where the binary is run. See the tipsandtechniques file locations section for more details. 

## vertica-prometheus-exporter global configuration settings
**Note the min_interval is set to 10 seconds. Collector files can have their own min_interval settings that override the global setting**
```yml
global:
# Scrape timeouts ensure that:
  #    (i)  scraping completes in reasonable time and
  #    (ii) slow queries are canceled early when the database is already under heavy load
  #  Prometheus informs targets of its own scrape timeout (via the "X-Prometheus-Scrape-Timeout-Seconds" request header)
  #  so the actual timeout is computed as:
  #    min(scrape_timeout, X-Prometheus-Scrape-Timeout-Seconds - scrape_timeout_offset)
  # 
  #  If scrape_timeout <= 0, no timeout is set unless Prometheus provides one. The default is 10s.
  scrape_timeout: 10s
  #  Subtracted from Prometheus' scrape_timeout to give us some headroom and prevent Prometheus from timing out first.
  # 
  #  Must be strictly positive. The default is 500ms.
  scrape_timeout_offset: 500ms
  #  Minimum interval between collector runs: by default (0s) collectors are executed on every scrape.
  min_interval: 10s
  #  Maximum number of open connections to any one target. Metric queries will run concurrently on multiple connections,
  #  as will concurrent scrapes.
   
  #  If max_connections <= 0, then there is no limit on the number of open connections. The default is 3.
  max_connections: 3
  #  Maximum number of idle connections to any one target. Unless you use very long collection intervals, this should
  #  always be the same as max_connections.
  
  #  If max_idle_connections <= 0, no idle connections are retained. The default is 3.
  max_idle_connections: 3
  #  Maximum number of maximum amount of time a connection may be reused. Expired connections may be closed lazily before reuse.
  #  If 0, connections are not closed due to a connection's age.
  max_connection_lifetime: 5m

# The target to monitor and the collectors to execute on it.
target:
  #  Data source name always has a URI schema that matches the driver name. In some cases (e.g. vertica)
  #  the schema gets dropped or replaced to match the driver expected DSN format.
  data_source_name: 'vertica://<username>:<userpwd>@<exporterhostip>:5433/<databasename>'

  #  Collectors (referenced by name) to execute on the target.
  collectors: [vertica_base_graphs ,vertica_base_gauges]

# Collector files specifies a list of globs. One collector definition is read from each matching file.
collector_files: 
- "*.collector.yml"
Log:
# Any integer value which represents days . 
  retention_day:  1 
# Any integer value which represents log file size in  megabytes 
  max_log_filesize:  1 

```

#### LOAD BALANCE AND FAIL SAFE: 
You can add Vertica native connection_load_balance and backup_server_node parameters via the data source name in the vertica-prometheus-exporter.yml file for best distributed connections and fail safety in event the primary Vertica node goes down. See the Vertica documentation for more details on these two features.

Here's an example of them added to the basic data_source_name string
```
  *data_source_name: 'vertica://dbadmin:@nn.nn.nn.235:5433/VMart?connection_load_balance=1&backup_server_node=nn.nn.nn.236:5433,nn.nn.nn.237:5433'*
```

#### TLS AUTHENTICATION: 
You can add Vertica TLS authentication parameter via the data source name in the vertica-prometheus-exporter.yml file. See the Vertica documentation for more details on TLS/SSL authenticaiton and server side setup requirements. The exporter supports tlsmodes of "none" and "server". It doesn't support "server-strict".

Here's an example of the tlsmode added to the basic data_source_name string
```
  *data_source_name: 'vertica://dbadmin:@nn.nn.nn.235:5433/VMart?tlsmode=server*
```



## vertica-base-example.collector.yml configuration
**Note in this file that all metrics issue their own query, and all queries return a single numeric value result.**
```yml
collector_name: example
### min_interval: 0s
metrics:
   - metric_name: vertica_license_size
     type: gauge
     help: 'Total License size in Bytes'
     values: [licsz]
     query: |
         select  /*+ LABEL(exporter_vertica_license_size MB) */ (license_size_bytes/1000000)::INTEGER as licsz
         from license_audits where audited_data='Total'
         order by audit_end_timestamp desc limit 1;
   - metric_name: vertica_database_size
     type: gauge
     help: 'Total Database size in MB'
     values: [ttldbsz]
     query: |
         select  /*+ LABEL(exporter_vertica_total_database_size) */ (database_size_bytes/1000000)::INTEGER as ttldbsz
         from license_audits where audited_data='Total'
         order by audit_end_timestamp desc limit 1;
   - metric_name: vertica_total_database_rows
     type: gauge
     help: 'Total Rows in Database from projection_storage table.'
     values: [ttlrows]
     query: |
         select /*+ LABEL(exporter_vertica_total_projection_rows) */ sum(row_count) as ttlrows 
         from projection_storage;
   - metric_name: vertica_total_database_connections
     type: gauge
     help: 'Total Database Connections from sessions table.'
     values: [ttlconns]
     query: |
         select /*+ LABEL(exporter_vertica_total_database_connections) */ count(*) as ttlconns 
         from sessions;
   - metric_name: vertica_state_not_up_or_standby
     type: gauge
     help: 'Nodes with state of other than UP or STANDBY.'
     values: [down]
     query: |
         select count(*) as down 
         from nodes 
         where node_state!='UP' and node_state!='Standby';

```

## vertica-base-example1.collector.yml configuration
**Note in this file that some metrics have their own queries and some reference different values returned by the single query at the end.**
```yml
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
  - metric_name: vertica_query_requests_transactions_count_per_node
    type: gauge
    help: 'Running transactions per node'
    key_labels:
       - node_name
    values: [total]
    query: |
       SELECT /*+ LABEL(exporter_vertica_query_requests_transactions_count_per_node) */
       node_name , count(*) total
       FROM transactions
       WHERE start_timestamp between date_trunc('minute',sysdate) - '1 minutes'::interval and date_trunc('minute',sysdate) - '1 milliseconds'::interval
       GROUP BY node_name
       ORDER BY node_name;
  - metric_name: vertica_cpu_usage_pct
    type: gauge
    help: 'vertica cpu usage percentage'
    key_labels: 
       - node_name
    values: [avg_cpu_usage_pct]
    query_ref: vertica_system_resources
  - metric_name: vertica_mem_usage_pct
    type: gauge
    help: 'vertica memory usage percentage'
    key_labels:
       - node_name
    values: [avg_mem_usage_pct]
    query_ref: vertica_system_resources
  - metric_name: vertica_net_rx_bytespersec
    type: gauge
    help: 'Vertica Network Receive bps'
    key_labels:
       - node_name
    values: [net_rx_bps]
    query_ref: vertica_system_resources
  - metric_name: vertica_net_tx_bytespersec
    type: gauge
    help: 'Vertica Network Transmit bps'
    key_labels:
       - node_name
    values: [net_tx_bps]
    query_ref: vertica_system_resources
  - metric_name: vertica_io_read_bytespersec
    type: gauge
    help: 'Vertica IO Read bps'
    key_labels:
       - node_name
    values: [io_read_bps]
    query_ref: vertica_system_resources
  - metric_name: vertica_io_write_bytespersec
    type: gauge
    help: 'Vertica IO Writes bps'
    key_labels:
       - node_name
    values: [io_write_bps]
    query_ref: vertica_system_resources

queries:
  - query_name: vertica_system_resources
    query: |
       select  
          node_name,
          ROUND(max(average_cpu_usage_percent)) as avg_cpu_usage_pct,
          ROUND(max(average_memory_usage_percent)) as avg_mem_usage_pct,
          CAST(max(net_rx_kbytes_per_second)*1024 as INTEGER) as net_rx_bps,
          CAST(max(net_tx_kbytes_per_second)*1024 as INTEGER) as net_tx_bps,
          CAST(max(io_read_kbytes_per_second)*1024 as INTEGER) as io_read_bps,
          CAST(max(io_written_kbytes_per_second)*1024 as INTEGER) as io_write_bps
       from system_resource_usage 
       group by node_name
       order by end_time desc limit 1;
       
```

### Configurations :

Global Configurations
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
   Maximum number of open connections to any one target. Metric queries will run concurrently on multiple connections,
    as will concurrent scrapes.
   
   If max_connections <= 0, then there is no limit on the number of open connections. The default is 3.
  max_connections: 3
   Maximum number of idle connections to any one target. Unless you use very long collection intervals, this should
   always be the same as max_connections.
  
   If max_idle_connections <= 0, no idle connections are retained. The default is 3.
  max_idle_connections: 3
  Maximum number of maximum amount of time a connection may be reused. Expired connections may be closed lazily before reuse.
  If 0, connections are not closed due to a connection's age.
  max_connection_lifetime: 5m

# The target to monitor and the collectors to execute on it.
target:
   Data source name always has a URI schema that matches the driver name. In some cases (e.g. vertica)
   the schema gets dropped or replaced to match the driver expected DSN format.
    data_source_name: 'vertica://<username>:<userpwd>@<exporterhostip>:5433/<databasename>'

  # Collectors (referenced by name) to execute on the target.
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
Vertica Base example Configurations
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

Vertica Base example1 Configurations
```yml
collector_name: vertica_example1
metrics:
  - metric_name: vertica_connections_per_node
    type: counter
    help: 'Connections per node'
    key_labels:
       - node_name
    values: [ttlconnspn]
    query: |
         select /*+ LABEL(exporter_vertica_total_database_connections) */ node_name,count(*) as ttlconnspn 
         from sessions
         group by node_name
         order by node_name;
  - metric_name: vertica_query_requests_transactions_count_per_node
    type: counter
    help: 'Running transactions per node'
    key_labels:
       - node_name
    values: [total]
    query: |
       SELECT /*+ LABEL(exporter_vertica_query_requests_transactions_count_per_node) */
       node_name , count(*) total
       FROM transactions
       WHERE start_timestamp between date_trunc('minute',sysdate) - '1 minutes'::interval and date_trunc('minute',sysdate) - '1 milliseconds'::interval
       GROUP BY node_name;
  - metric_name: vertica_cpu_usage_pct
    type: counter
    help: 'vertica cpu usage percentage'
    key_labels: 
       - node_name
    values: [avg_cpu_usage_pct]
    query_ref: vertica_system_resources
  - metric_name: vertica_mem_usage_pct
    type: counter
    help: 'vertica memory usage percentage'
    key_labels:
       - node_name
    values: [avg_mem_usage_pct]
    query_ref: vertica_system_resources
  - metric_name: vertica_net_rx_kbytespersec
    type: counter
    help: 'Vertica Network Receive kbps'
    key_labels:
       - node_name
    values: [net_rx_kbps]
    query_ref: vertica_system_resources
  - metric_name: vertica_net_tx_kbytespersec
    type: counter
    help: 'Vertica Network Transmit kbps'
    key_labels:
       - node_name
    values: [net_tx_kbps]
    query_ref: vertica_system_resources
  - metric_name: vertica_io_read_kbytespersec
    type: counter
    help: 'Vertica IO Read kbps'
    key_labels:
       - node_name
    values: [io_read_kbps]
    query_ref: vertica_system_resources
  - metric_name: vertica_io_write_kbytespersec
    type: counter
    help: 'Vertica IO Writes kbps'
    key_labels:
       - node_name
    values: [io_write_kbps]
    query_ref: vertica_system_resources

queries:
  - query_name: vertica_system_resources
    query: |
       select  
          node_name,
          max(average_cpu_usage_percent) as avg_cpu_usage_pct,
          max(average_memory_usage_percent) as avg_mem_usage_pct,
          max(net_rx_kbytes_per_second) as net_rx_kbps,
          max(net_tx_kbytes_per_second) as net_tx_kbps,
          max(io_read_kbytes_per_second) as io_read_kbps,
          max(io_
       - node_name
    values: [io_write_kbps]
    query_ref: vertica_system_resources

```

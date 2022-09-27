### WHEN TO USE: 

Comparison of available tools – 

- Vertica Management Console – Comes with Vertica but requires seperate install. All metrics available are chosen by Vertica for value to user. It's a static set so user can’t customize the metrics available or how displayed.  

- Grafana plugin – Requires Grafana with Vertica datasource installed and configured. The data source does a direct connection to Vertica to get metrics. There is no  persistent retention and it runs queries based on the defined Grafana panel interval. Queries and dashboards fully customizable so user chooses what they want to look at and how. Supports tables with string values. Requires Vertica SQL and Grafana skills  

- Prometheus exporter – Requires Prometheus and vertica-prometheus-exporter. Prometheus is configured to use the vertica-prometheus-exporter as a target and initiates scrapes from it. The exporter connects to Vertica, runs the metrics queries, and returns them to Prometheus in a known format. The exporter has the ability to cache metrics if desired, which can help minimize the scrapes it does to the Vertica database. The Prometheus metrics can be used by any tool that supports Prometheus format, e.g. Grafana. Collectors (sets of metrics queries) are fully customizable so user chooses what they want to look at. Requires Vertica SQL and Prometheus skills, plus whatever evisualization tool is being used skills.  

### LOAD BALANCE AND FAIL SAFE: 

Activate connection_load_balance  and backup_server_node  via the data source name in the vertica_prometheus_exporter.yml file for best distributed connections and fail safety in event the primary Vertica node goes down. See the Vertica documentaiton for more details on these ttwo features.

  *data_source_name: 'vertica://dbadmin:@nn.nn.nn.235:5433/VMart?connection_load_balance=1&backup_server_node=nn.nn.nn.236:5433,nn.nn.nn.237:5433'*
  

### STARTUP ORDER: 

- First start the vertica-prometheus-exporter. Depending on your deployment It should say listening at the end of console output, end of Logfile/logfile.log, or end of nohup.out. 

- Wait a minute and then start Prometheus.  It should say listening at the end of console output or end of nohup.out. 

- Now go to the Prometheus http interface (http://<prometheusserverip>:9090/targets. If the status of the vertica-export says Down or Unknown wait 30 seconds or so then refresh.. Repeat until it says Up. Now you can view the metrics  

### NAMES: 

Keep your metric names for metrics where node_name is a label short. The combination of the metric name plus the long Vertica node path can result in truncation in the Grafana panels or force you to make them wider than planned. 

**TYPE vertica_query_requests_transactions_count_per_node counter** 
```
vertica_query_requests_transactions_count_per_node{node_name="v_vmart_node0001"} 11 
vertica_query_requests_transactions_count_per_node{node_name="v_vmart_node0002"} 15 
vertica_query_requests_transactions_count_per_node{node_name="v_vmart_node0003"} 9 
```

### FILE LOCATIONS 

**Non Docker Linux Build** 

Once you build the vertica_prometheus_exporter binary you can move it to any location you want but the following dependencies exist: 
- You have to launch the binary from the directory where it exists 
- You have to have a metrices dir under the binary’s directory that contains the collector yml files 
- You can have the vertica_prometheus_exporter.yml file anywhere you like as the –config-file parameter to starting the binary can be a fully qualified path to it. 
- Your Logfile dir must be under the binary’s directory. The binary will create the Logfile dir and the logfile.log in it if they don’t exist. 
 

**Docker build** 

To be added. Note same as above, but possible to use -v to bind dirs. external to container 

### MINIMIZE QUERY IMPACT on VERTICA 

There is interval control at several levels (end tool such as Grafana, Prometheus, and the exporter). Make sure to set the intervals for the best efficiency. Don’t collect slow changing values frequently. Maybe group collections by rate of change and frequency of scrape. See the min_intervals section for more details. 

You can set the exporter max_connections to set how many concurrent connections a metrics scrape will establish. There is no pooling, all connections end as soon as their task is complete. The number of connections will dictate the duration of the metric scrape, but could impact normal operations if set too high. Find a balance between time and resources. 

If planning on rendering several metrics from the same table use a query object asking for all columns desired and then queryrefs to get those columns into individual panels. This will issue one query vs several each interval refresh. 

### PROMETHEUS STORAGE 

Prometheus by default retains 15 days of metric data. See the Prometheus documentation for location of the database, how to adjust the retention values for size and/or time, and best practices. This is just to note to raise awareness that the more metrics you have the vertica prometheus exporter capture the larger the Prometheus database will become. It should be monitored and space issues resolved. 


### EXPORTER STORAGE 

To exporter has a relatively small footprint, but it does produce a logfile that can grow over time. In the vertica_prometheus_exporter.yml file we’ve provided two knobs to manage how large the log file will get. There are knobs to define number of days retained and max file size. 
```yml
Log: 
  retention_day:  15 
  max_log_filesize:  500 # in megabytes 
```

### CLEARTEXT PASSWORD 

To prevent passing the database password in cleartext across the lan we’ve used a secret to store it. The password will still be clear text in the vertica_prometheus_exporter.yml so that file should have an access mask to only allow the user running the exporter to read it.  

### COLLECTOR FILE PLACEMENT 

To prevent false console/log output by the Exporter, only put collector yml files in the metrices directory that are associated with collectors you specify in the vertica_prometheus_exporter.yml config file. Alternatively, if you want to keep a superset of collectors but change which you use at different times, then instead of the collector_files value specifying a glob (*.collector.yml) , list the yml files individually (vertica_base_graphs,vertica_base_gauges)..  

If yml files not associated with the collectors: value exist in the collector files dir the console output will imply all collectors were loaded. 

### Example: 

If we have vertica_base_gauges.yml and vertica_base_graphs.yml in my yml dir, but in my vertica_prometheus_exporter.yml config file we only specify the graphs and we use the glob, the exporter start output says it loaded the gauges and graphs. You can verify in Prometheus that it in fact only loaded the graphs as specified. 

Vertica_export.yml extract 

**Collectors (referenced by name) to execute on the target.** 
```  
  collectors: [vertica_base_graphs] 
```
**Collector files specifies a list of globs. One collector definition is read from each matching file.**
```
collector_files: 
  - "*.collector.yml" 
```
Exporter startup output 

`./vertica-prometheus-exporter --config.file vertica.yml` 
```
I0902 11:35:26.047056  140624 main.go:63] Starting vertica prometheus exporter (version=, branch=, revision=) (go=go1.18.4, user=, date=) 
I0902 11:35:26.047267  140624 config.go:22] Loading configuration from vertica.yml 
I0902 11:35:26.048940  140624 config.go:148] Loaded collector "vertica_base_gauges" from vertica_base_gauges.collector.yml 
I0902 11:35:26.049597  140624 config.go:148] Loaded collector "vertica_base_graphs" from vertica_base_graphs.collector.yml 
(0xa4df40,0xc000194c80) 
I0902 11:35:26.049919  140624 main.go:82] Listening on :9968 
```
### vertica prometheus exporter min_interval EXPLAINED 

The min_interval knob determines the lifespan of the internal collector objects. A collector with min_interval=0s will open, scrape Vertica, and close. A collector with min_interval=60s will open, scrape Vertica, and remain open as a temporary cache. Subsequent requests for that collector from Prometheus prior to the min_interval will get cached results from the exporter and not scrape Vertica. A request for that collector from Prometheus which occurs after the min_interval is reached will get a new collector, fresh scrape of Vertica, and again live for the duration of min_interval. 

There is a global min_interval setting in the vertica.yml file. This governs the min_interval for all active collectors. Each collector file can have it’s own min_interval setting, allowing you to control how frequently a Prometheus request actually causes a scrape against Vertica. 

`[dbadmin@vertica-node metrices]$ head vertica_prometheus_exporter.yml `
```
global: 
  scrape_timeout_offset: 500ms 
  min_interval: 10s 
```
`[dbadmin@vertica-node metrices]$ head vertica_base_graphs.collector.yml`
```
collector_name: vertica_base_graphs 
min_interval: 75s 
```
`[dbadmin@vertica-node metrices]$ head vertica_base_gauges.collector.yml`
```
collector_name: vertica_base_gauges 
min_interval: 0s 
```
In the exporter log file you can see the collectors being created, aged, and the return of cached metrics vs fresh metrics. Above you can see that global value is 10s, and graphs collection is set to 75s. Below the log entries for Prometheus requests for graph metrics show the new collector, return of cached metrics, age reached, and return of fresh metrics. 

**first collection collector age is a large number, subsequent ones will return normal values for age:** 
```
INFO[2022-09-12T15:08:51-04:00] [collector="vertica_base_graphs"] Collecting fresh metrics: min_interval=75.000s cache_age=9223372036.855s 
```
**here we see that Prometheus will get cached metrics because the collector has not yet reached min_interval (lifespan)**
```
INFO[2022-09-12T15:09:51-04:00] [collector="vertica_base_graphs"] Returning cached metrics: min_interval=75.000s cache_age=59.990s 
```
**when the cache reaches min_interval (lifespan) then Prometheus will get a new collector and new metrics from Vertica:**
```
INFO[2022-09-12T15:10:51-04:00] [collector="vertica_base_graphs"] Collecting fresh metrics: min_interval=75.000s cache_age=119.994s 
```
**you can see the difference between cached and fresh. The fresh shows columns returned from Vertica** 
```
[2022-09-12T15:09:51-04:00] [collector="vertica_base_graphs"] Returning cached metrics: min_interval=75.000s cache_age=59.990s 

INFO[2022-09-12T15:10:51-04:00] [collector="vertica_base_graphs"] Collecting fresh metrics: min_interval=75.000s cache_age=119.994s 
INFO[2022-09-12T15:10:51-04:00] returned_columns="[node_name totaltrans]"collector="vertica_base_graphs", query="vertica_connections_per_node" 
INFO[2022-09-12T15:10:51-04:00] returned_columns="[node_name total]"collector="vertica_base_graphs", query="vertica_query_requests_transactions_count_per_node" 
INFO[2022-09-12T15:10:54-04:00] returned_columns="[node_name avg_cpu_usage_pct avg_mem_usage_pct net_rx_bps net_tx_bps io_read_bps io_write_bps]"collector="vertica_base_graphs", query="vertica_system_resources" 
```

**with my collector min_interval=0s every Prometheus request gets fresh metrics** 
```
INFO[2022-09-12T15:36:30-04:00] Listening on :9968 
INFO[2022-09-12T15:36:51-04:00] returned_columns="[node_name totaltrans]"collector="vertica_base_graphs", query="vertica_connections_per_node" 
INFO[2022-09-12T15:36:51-04:00] returned_columns="[node_name total]"collector="vertica_base_graphs", query="vertica_query_requests_transactions_count_per_node" 
INFO[2022-09-12T15:36:54-04:00] returned_columns="[node_name avg_cpu_usage_pct avg_mem_usage_pct net_rx_bps net_tx_bps io_read_bps io_write_bps]"collector="vertica_base_graphs", query="vertica_system_resources" 
INFO[2022-09-12T15:37:51-04:00] returned_columns="[node_name totaltrans]"collector="vertica_base_graphs", query="vertica_connections_per_node" 
INFO[2022-09-12T15:37:51-04:00] returned_columns="[node_name total]"collector="vertica_base_graphs", query="vertica_query_requests_transactions_count_per_node" 
INFO[2022-09-12T15:37:54-04:00] returned_columns="[node_name avg_cpu_usage_pct avg_mem_usage_pct net_rx_bps net_tx_bps io_read_bps io_write_bps]"collector="vertica_base_graphs", query="vertica_system_resources" 
```
So using the min_intervals at the collector level allows you to control how often Prometheus gets up to date results, and how often Vertica is queried. You could potentially create collectors which you want real time results every request (min_interval=0s), reasonably real time results for fast moving metrics (min_interval=60s), and periodic results for slow moving metrics (min_interval=600s). 

*Note: if you set min_interval shorter than the Prometheus scrape interval you will get a fresh scrape every time. So make sure to review Prometheus scrape interval before setting the exporter min_intervals.* 

*Tuning hint. If you are using min_interval > 0s you can grep the log file to compare the number of cached (no Vertica scrape) to fresh (Vertica scrape) results. This may help you decide if you want to adjust for more or fewer Vertica scrapes for the particular metric.* 

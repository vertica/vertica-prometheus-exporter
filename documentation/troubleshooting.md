## TROUBLESHOOTING
This document covers some common errors one might see and the potential root causes.

## UNUSED COLUMNS IN QUERY
If you see these in the exporter output they are warnings.
```
W0817 16:15:51.296708   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "end_time" returned by query
W0817 16:15:51.296759   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "s3ops" returned by query
W0817 16:15:51.296778   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "s3errs" returned by query
W0817 16:15:51.296797   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "s3retries" returned by query
W0817 16:15:51.296815   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "s3listheads" returned by query
W0817 16:15:51.296832   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "s3postdeletes" returned by query
W0817 16:15:51.296850   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "bytessent" returned by query
W0817 16:15:51.296868   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_s3_performance"] Extra column "bytesrecvd" returned by query
W0817 16:15:51.305122   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_depot_evictions_by_minute"] Extra column "YMDHM" returned by query
W0817 16:15:51.306289   96868 query.go:150] [collector="vertica_base_EON_Mode", query="vertica_depot_fetches_by_minute"] Extra column "YMDHM" returned by query
```
Root cause is typically due to a query object being defined and not all columns in the query being used in queryref objects. E.g. the query below only has one metric using it in a single value queryref. So all the other columns in the query are flagged as extras not used. To fix this one would change the query object to only select columns that will have metrics that queryref them.

```yml
  - metric_name: vertica_data_writes_per_hour
    type: counter
    help: 'S3 Data Writes (puts)'
    key_labels:
       - node_name
    values: [s3puts]
    query_ref: vertica_s3_performance

  - query_name: vertica_s3_performance
    query: |
        select node_name,end_time,avg_operations_per_second as s3ops, avg_errors_per_second as s3errs, retries as s3retries,
        metadata_reads as s3listheads, metadata_writes as s3postdeletes, data_reads as s3gets, data_writes as s3puts,
        upstream_bytes as bytessent, downstream_bytes as bytesrecvd
        from udfs_ops_per_hour
        where filesystem='S3'
        order by node_name,end_time desc;
```

## TABLES WITH STRING VALUES
Prometheus doesn't support string metrics, so any attempt at creating a metric with non numeric return values will error when trying to convert.
https://stackoverflow.com/questions/65850083/prometheus-java-client-export-string-based-metrics
https://github.com/prometheus/prometheus/issues/2227

Return values are integer:

```
dbadmin=> \d mvaltesti
                                  List of Fields by Tables
 Schema |   Table   | Column | Type | Size | Default | Not Null | Primary Key | Foreign Key
--------+-----------+--------+------+------+---------+----------+-------------+-------------
 public | mvaltesti | c1     | int  |    8 |         | f        | f           |
 public | mvaltesti | c2     | int  |    8 |         | f        | f           |
 public | mvaltesti | c3     | int  |    8 |         | f        | f           |
```

collector_name: vertica_base_tables
metrics:
  - metric_name: vertica_node_states
    type: gauge
    help: 'State of all Nodes in the Database'
    key_labels:
       - c1
    value_label: colstoget
    values:
        - c2
        - c3
    query: |
        select c1,c2,c3
        from mvaltesti;

```
INFO[2022-09-19T13:48:51-04:00] [collector="vertica_base_tables"] Collecting fresh metrics: min_interval=10.000s cache_age=60.012s
INFO[2022-09-19T13:48:51-04:00] returned_columns="[c1 c2 c3]"collector="vertica_base_tables", query="vertica_node_states"
```

Return values are Non integer:

```
dbadmin=> \d mvaltest
                                     List of Fields by Tables
 Schema |  Table   | Column |    Type     | Size | Default | Not Null | Primary Key | Foreign Key
--------+----------+--------+-------------+------+---------+----------+-------------+-------------
 public | mvaltest | c1     | varchar(10) |   10 |         | f        | f           |
 public | mvaltest | c2     | varchar(10) |   10 |         | f        | f           |
 public | mvaltest | c3     | varchar(10) |   10 |         | f        | f           |
```

collector_name: vertica_base_tables
metrics:
  - metric_name: vertica_node_states
    type: gauge
    help: 'State of all Nodes in the Database'
    key_labels:
       - c1
    value_label: colstoget
    values:
        - c2
        - c3
    query: |
        select c1,c2,c3
        from mvaltest;
```
INFO[2022-09-19T13:53:51-04:00] [collector="vertica_base_tables"] Collecting fresh metrics: min_interval=10.000s cache_age=9223372036.855s
INFO[2022-09-19T13:53:51-04:00] returned_columns="[c1 c2 c3]"collector="vertica_base_tables", query="vertica_node_states"
INFO[2022-09-19T13:53:51-04:00] Error gathering metrics:%!(EXTRA prometheus.MultiError=2 error(s) occurred:
* [from Gatherer #1] [collector="vertica_base_tables", query="vertica_node_states"] scanning of query result failed: sql: Scan error on column index 1, name "c2": converting driver.Value type string ("two") to a float64: invalid syntax
* [from Gatherer #1] [collector="vertica_base_tables", query="vertica_node_states"] scanning of query result failed: sql: Scan error on column index 1, name "c2": converting driver.Value type string ("two2") to a float64: invalid syntax)
```

### EMPTY OR NULL COLUMN VALUES
Any column returning an empty or null value will fail. A system table example would be the node_down_since columns of the nodes table. That column is empty unless the node state is DOWN.
```
dbadmin=> select * from mvaltesti;
 c1 | c2 | c3
----+----+----
 11 | 22 | 33
  1 |  2 |
```
```
INFO[2022-09-19T13:58:51-04:00] [collector="vertica_base_tables"] Collecting fresh metrics: min_interval=10.000s cache_age=9223372036.855s
INFO[2022-09-19T13:58:51-04:00] returned_columns="[c1 c2 c3]"collector="vertica_base_tables", query="vertica_node_states"
INFO[2022-09-19T13:58:51-04:00] Error gathering metrics:%!(EXTRA *errors.errorString=[from Gatherer #1] [collector="vertica_base_tables", query="vertica_node_states"] scanning of query result failed: sql: Scan error on column index 2, name "c3": converting NULL to float64 is unsupported)
```

## DATA GAPS IN GRAPHS
Prometheus and tools like Grafana may show gaps in the data in graphs for some metrics. This is most likely because there is no data for one or more rows for that metric scrape at that time. It will happen more often on an idle Vertica database than a busy one

Here's an example of a metric against a 3 node database where at the time of the scrape only node0001 had activity
```  
metric_name: vertica_query_requests_transactions_count_per_node
SELECT /*+ LABEL(exporter_vertica_query_requests_transactions_count_per_node) */ node_name , count(*) total FROM transactions 
WHERE start_timestamp between date_trunc('minute',sysdate) - '1 minutes'::interval and date_trunc('minute',sysdate) - 
'1 milliseconds'::interval GROUP BY node_name;
```
vsql shows that there are only 4 transactions on node0001, nothing for node0002 or node0003
``` 
    node_name     | total
------------------+-------
 v_vmart_node0001 |     4
(1 row)
```
Prometheus metrics output shows only a row for the node0001 that had a return value, no rows for node0002 or node0003
```
# HELP vertica_query_requests_transactions_count_per_node Running transactions per node
# TYPE vertica_query_requests_transactions_count_per_node gauge
vertica_query_requests_transactions_count_per_node{node_name="v_vmart_node0001"} 4
```
In Prometheus graphs, or tools like Grafana, this will show the missing values as gaps in the graph for the line/bar representing the nodes that had no return values. Prometheus has a Resolution setting you can adjust that may help adjust the graph resolution to minimize the gaps visually. Grafana has a "Connect null values" option on the graph panels that will fill the gaps, noting that this could be technically misleading even though visually pleasing. 

If your Vertica query happens to be a timeseries you may be able to use the TS_LATEST_VALUE or TS_FIRST_VALUE functions to fill gaps if desired. See the Vertica documentation for more details.

## TLSMODE SERVER-STRICT FAILS
Currently the exporter supports using the data_source_name tlsmode parameters of either tlsmode=none or tlsmode=server.

tlsmode server-strict is not currently implemented

If you try to use strict (mutual) mode you will get an error like this in the exporter log or console
Oct  7 12:27:24.542147 ERROR driver: x509: certificate signed by unknown authority
INFO[2022-10-07T12:27:24-04:00] Error gathering metrics:%!(EXTRA *errors.errorString=[from Gatherer #1] x509: certificate signed by unknown authority)

If you turn on the vertica-sql-go driver debug logging you will see a similar message in it's log file.
go drive debug level log shows the same
Oct  7 13:12:41.702736 ERROR connection: -> FAILED SENDING Startup (packet): ProtocolVersion:00030009, DriverName='vertica-sql-go', DriverVersion='1.2.2', UserName='dbadmin', Database='VMart', SessionID='vertica-sql-go-1.2.2-257565-1665162761', ClientPID=257565: x509: certificate signed by unknown authority
Oct  7 13:12:41.702814 ERROR driver: x509: certificate signed by unknown authority



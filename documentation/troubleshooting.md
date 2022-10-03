
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
Root cause is typically due to a query object being defined and not all columns in the query being used in queryref objects. E.g.the query below only has one metric using it in a single value queryref. So all the other columns in the query are flagged as extras not used. To fix this one would change the query object to only select columns that will have metrics that queryref them.

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
Prometheus doesn't support string metrics, so any attempt at creating a metric with non numeric return vlaues will error when trying to convert.
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

## EMPTY OR NULL VALUES
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




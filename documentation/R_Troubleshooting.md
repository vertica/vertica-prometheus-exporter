--------------------
Unused Columns
--------------------
If you see these in the exporter output they are warnings.
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

Root cause is typically due to a query object being defined and not all columns in the query being used in queryref objects. E.g.the query below only ha sone metric using it in a single value queryref. So all the other columns in the query are flagged as extras not used. To fix this one would change the query object to only select columns that will have metrics that queryrefthem.

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

--------------------

--------------------
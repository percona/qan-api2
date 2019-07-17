CREATE TABLE metrics (
  -- Main dimensions
  `queryid` LowCardinality(String) COMMENT 'hash of query fingerprint',
  `server` LowCardinality(String) COMMENT 'IP or hostname of DB server',
  `database` LowCardinality(String) COMMENT 'PostgreSQL: database',
  `schema` LowCardinality(String) COMMENT 'MySQL: database; PostgreSQL: schema',
  `username` LowCardinality(String) COMMENT 'client user name',
  `client_host` LowCardinality(String) COMMENT 'client IP or hostname',
  --  Standard labels
  `replication_set` LowCardinality(String) COMMENT 'Name of replication set',
  `cluster` LowCardinality(String) COMMENT 'Cluster name',
  `service_type` LowCardinality(String) COMMENT 'Type of service',
  `environment` LowCardinality(String) COMMENT 'Environment name',
  `az` LowCardinality(String) COMMENT 'Availability zone',
  `region` LowCardinality(String) COMMENT 'Region name',
  `node_model` LowCardinality(String) COMMENT 'Node model',
  `container_name` LowCardinality(String) COMMENT 'Container name',
  -- Custom labels
  `labels.key` Array(LowCardinality(String)) COMMENT 'Custom labels names',
  `labels.value` Array(LowCardinality(String)) COMMENT 'Custom labels values',
  `agent_id` LowCardinality(String) COMMENT 'Identifier of agent that collect and send metrics',
  `agent_type` Enum8(
    'agent_type_invalid' = 0,
    'mysql-perfschema' = 1,
    'mysql-slowlog' = 2,
    'mongodb-profiler' = 3
  ) COMMENT 'Agent Type that collect of metrics: slowlog, perf schema, etc.',
  `period_start` DateTime COMMENT 'Time when collection of bucket started',
  `period_length` UInt32 COMMENT 'Duration of collection bucket',
  `fingerprint` LowCardinality(String) COMMENT 'mysql digest_text; query without data',
  `example` String COMMENT 'One of query example from set found in bucket',
  `example_format` Enum8(
    'EXAMPLE_FORMAT_INVALID' = 0,
    'EXAMPLE' = 1,
    'FINGERPRINT' = 2
  ) COMMENT 'Indicates that collect real query examples is prohibited',
  `is_truncated` UInt8 COMMENT 'Indicates if query examples is too long and was truncated',
  `example_type` Enum8(
    'EXAMPLE_TYPE_INVALID' = 0,
    'RANDOM' = 1,
    'SLOWEST' = 2,
    'FASTEST' = 3,
    'WITH_ERROR' = 4
  ) COMMENT 'Indicates what query example was picked up',
  `example_metrics` String COMMENT 'Metrics of query example in JSON format.',
  `num_queries_with_warnings` Float32 COMMENT 'How many queries was with warnings in bucket',
  `warnings.code` Array(UInt32) COMMENT 'List of warnings',
  `warnings.count` Array(Float32) COMMENT 'Count of each warnings in bucket',
  `num_queries_with_errors` Float32 COMMENT 'How many queries was with error in bucket',
  `errors.code` Array(UInt64) COMMENT 'List of Last_errno',
  `errors.count` Array(UInt64) COMMENT 'Count of each Last_errno in bucket',
  `num_queries` Float32 COMMENT 'Amount queries in this bucket',
  -- Metrics
  `m_query_time_cnt` Float32 COMMENT 'The statement execution time in seconds was met.',
  `m_query_time_sum` Float32 COMMENT 'The statement execution time in seconds.',
  `m_query_time_min` Float32 COMMENT 'Smallest value of query_time in bucket',
  `m_query_time_max` Float32 COMMENT 'Biggest value of query_time in bucket',
  `m_query_time_p99` Float32 COMMENT '99 percentile of value of query_time in bucket',
  `m_lock_time_cnt` Float32,
  `m_lock_time_sum` Float32 COMMENT 'The time to acquire locks in seconds.',
  `m_lock_time_min` Float32,
  `m_lock_time_max` Float32,
  `m_lock_time_p99` Float32,
  `m_rows_sent_cnt` Float32,
  `m_rows_sent_sum` Float32 COMMENT 'The number of rows sent to the client.',
  `m_rows_sent_min` Float32,
  `m_rows_sent_max` Float32,
  `m_rows_sent_p99` Float32,
  `m_rows_examined_cnt` Float32,
  `m_rows_examined_sum` Float32 COMMENT 'Number of rows scanned - SELECT.',
  `m_rows_examined_min` Float32,
  `m_rows_examined_max` Float32,
  `m_rows_examined_p99` Float32,
  `m_rows_affected_cnt` Float32,
  `m_rows_affected_sum` Float32 COMMENT 'Number of rows changed - UPDATE, DELETE, INSERT.',
  `m_rows_affected_min` Float32,
  `m_rows_affected_max` Float32,
  `m_rows_affected_p99` Float32,
  `m_rows_read_cnt` Float32,
  `m_rows_read_sum` Float32 COMMENT 'The number of rows read from tables.',
  `m_rows_read_min` Float32,
  `m_rows_read_max` Float32,
  `m_rows_read_p99` Float32,
  `m_merge_passes_cnt` Float32,
  `m_merge_passes_sum` Float32 COMMENT 'The number of merge passes that the sort algorithm has had to do.',
  `m_merge_passes_min` Float32,
  `m_merge_passes_max` Float32,
  `m_merge_passes_p99` Float32,
  `m_innodb_io_r_ops_cnt` Float32,
  `m_innodb_io_r_ops_sum` Float32 COMMENT 'Counts the number of page read operations scheduled.',
  `m_innodb_io_r_ops_min` Float32,
  `m_innodb_io_r_ops_max` Float32,
  `m_innodb_io_r_ops_p99` Float32,
  `m_innodb_io_r_bytes_cnt` Float32,
  `m_innodb_io_r_bytes_sum` Float32 COMMENT 'Similar to innodb_IO_r_ops, but the unit is bytes.',
  `m_innodb_io_r_bytes_min` Float32,
  `m_innodb_io_r_bytes_max` Float32,
  `m_innodb_io_r_bytes_p99` Float32,
  `m_innodb_io_r_wait_cnt` Float32,
  `m_innodb_io_r_wait_sum` Float32 COMMENT 'Shows how long (in seconds) it took InnoDB to actually read the data from storage.',
  `m_innodb_io_r_wait_min` Float32,
  `m_innodb_io_r_wait_max` Float32,
  `m_innodb_io_r_wait_p99` Float32,
  `m_innodb_rec_lock_wait_cnt` Float32,
  `m_innodb_rec_lock_wait_sum` Float32 COMMENT 'Shows how long (in seconds) the query waited for row locks.',
  `m_innodb_rec_lock_wait_min` Float32,
  `m_innodb_rec_lock_wait_max` Float32,
  `m_innodb_rec_lock_wait_p99` Float32,
  `m_innodb_queue_wait_cnt` Float32,
  `m_innodb_queue_wait_sum` Float32 COMMENT 'Shows how long (in seconds) the query spent either waiting to enter the InnoDB queue or inside that queue waiting for execution.',
  `m_innodb_queue_wait_min` Float32,
  `m_innodb_queue_wait_max` Float32,
  `m_innodb_queue_wait_p99` Float32,
  `m_innodb_pages_distinct_cnt` Float32,
  `m_innodb_pages_distinct_sum` Float32 COMMENT 'Counts approximately the number of unique pages the query accessed.',
  `m_innodb_pages_distinct_min` Float32,
  `m_innodb_pages_distinct_max` Float32,
  `m_innodb_pages_distinct_p99` Float32,
  `m_query_length_cnt` Float32,
  `m_query_length_sum` Float32 COMMENT 'Shows how long the query is.',
  `m_query_length_min` Float32,
  `m_query_length_max` Float32,
  `m_query_length_p99` Float32,
  `m_bytes_sent_cnt` Float32,
  `m_bytes_sent_sum` Float32 COMMENT 'The number of bytes sent to all clients.',
  `m_bytes_sent_min` Float32,
  `m_bytes_sent_max` Float32,
  `m_bytes_sent_p99` Float32,
  `m_tmp_tables_cnt` Float32,
  `m_tmp_tables_sum` Float32 COMMENT 'Number of temporary tables created on memory for the query.',
  `m_tmp_tables_min` Float32,
  `m_tmp_tables_max` Float32,
  `m_tmp_tables_p99` Float32,
  `m_tmp_disk_tables_cnt` Float32,
  `m_tmp_disk_tables_sum` Float32 COMMENT 'Number of temporary tables created on disk for the query.',
  `m_tmp_disk_tables_min` Float32,
  `m_tmp_disk_tables_max` Float32,
  `m_tmp_disk_tables_p99` Float32,
  `m_tmp_table_sizes_cnt` Float32,
  `m_tmp_table_sizes_sum` Float32 COMMENT 'Total Size in bytes for all temporary tables used in the query.',
  `m_tmp_table_sizes_min` Float32,
  `m_tmp_table_sizes_max` Float32,
  `m_tmp_table_sizes_p99` Float32,
  -- Boolean metrics
  `m_qc_hit_cnt` Float32,
  `m_qc_hit_sum` Float32 COMMENT 'Query Cache hits.',
  `m_full_scan_cnt` Float32,
  `m_full_scan_sum` Float32 COMMENT 'The query performed a full table scan.',
  `m_full_join_cnt` Float32,
  `m_full_join_sum` Float32 COMMENT 'The query performed a full join (a join without indexes).',
  `m_tmp_table_cnt` Float32,
  `m_tmp_table_sum` Float32 COMMENT 'The query created an implicit internal temporary table.',
  `m_tmp_table_on_disk_cnt` Float32,
  `m_tmp_table_on_disk_sum` Float32 COMMENT 'The querys temporary table was stored on disk.',
  `m_filesort_cnt` Float32,
  `m_filesort_sum` Float32 COMMENT 'The query used a filesort.',
  `m_filesort_on_disk_cnt` Float32,
  `m_filesort_on_disk_sum` Float32 COMMENT 'The filesort was performed on disk.',
  `m_select_full_range_join_cnt` Float32,
  `m_select_full_range_join_sum` Float32 COMMENT 'The number of joins that used a range search on a reference table.',
  `m_select_range_cnt` Float32,
  `m_select_range_sum` Float32 COMMENT 'The number of joins that used ranges on the first table.',
  `m_select_range_check_cnt` Float32,
  `m_select_range_check_sum` Float32 COMMENT 'The number of joins without keys that check for key usage after each row.',
  `m_sort_range_cnt` Float32,
  `m_sort_range_sum` Float32 COMMENT 'The number of sorts that were done using ranges.',
  `m_sort_rows_cnt` Float32,
  `m_sort_rows_sum` Float32 COMMENT 'The number of sorted rows.',
  `m_sort_scan_cnt` Float32,
  `m_sort_scan_sum` Float32 COMMENT 'The number of sorts that were done by scanning the table.',
  `m_no_index_used_cnt` Float32,
  `m_no_index_used_sum` Float32 COMMENT 'The number of queries without index.',
  `m_no_good_index_used_cnt` Float32,
  `m_no_good_index_used_sum` Float32 COMMENT 'The number of queries without good index.',
  -- mongo metrics
  `m_docs_returned_cnt` Float32,
  `m_docs_returned_sum` Float32 COMMENT 'The number of returned documents.',
  `m_docs_returned_min` Float32,
  `m_docs_returned_max` Float32,
  `m_docs_returned_p99` Float32,
  `m_response_length_cnt` Float32,
  `m_response_length_sum` Float32 COMMENT 'The response length of the query result in bytes.',
  `m_response_length_min` Float32,
  `m_response_length_max` Float32,
  `m_response_length_p99` Float32,
  `m_docs_scanned_cnt` Float32,
  `m_docs_scanned_sum` Float32 COMMENT 'The number of scanned documents.',
  `m_docs_scanned_min` Float32,
  `m_docs_scanned_max` Float32,
  `m_docs_scanned_p99` Float32
) ENGINE = MergeTree PARTITION BY toYYYYMMDD(period_start)
ORDER BY
  (
    queryid,
    server,
    database,
    schema,
    username,
    client_host,
    period_start
  ) SETTINGS index_granularity = 8192;

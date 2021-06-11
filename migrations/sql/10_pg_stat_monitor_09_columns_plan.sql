ALTER TABLE metrics
  ADD COLUMN `m_wal_bytes` String COMMENT 'Bytes of WAL (Write-ahead logging) records',
  ADD COLUMN `m_plan_time_sum` Float32 COMMENT 'Sum of plan time.',
  ADD COLUMN `m_plan_time_min` Float32 COMMENT 'Min of plan time.',
  ADD COLUMN `m_plan_time_max` Float32 COMMENT 'Max of plan time.';

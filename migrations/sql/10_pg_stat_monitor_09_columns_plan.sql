ALTER TABLE metrics
  ADD COLUMN `m_wal_bytes` String COMMENT 'Bytes of WAL (Write-ahead logging) records',
  ADD COLUMN `m_plan_total_time` Float32 COMMENT 'Sum of plan time.',
  ADD COLUMN `m_plan_min_time` Float32 COMMENT 'Min of plan time.',
  ADD COLUMN `m_plan_max_time` Float32 COMMENT 'Max of plan time.',
  ADD COLUMN `m_plan_mean_time` Float32 COMMENT 'Mean of plan time.';

ALTER TABLE metrics
  DROP COLUMN `m_wal_bytes`,
  DROP COLUMN `m_plan_total_time`,
  DROP COLUMN `m_plan_min_time`,
  DROP COLUMN `m_plan_max_time`,
  DROP COLUMN `m_plan_mean_time`;

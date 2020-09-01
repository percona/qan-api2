ALTER TABLE metrics
  DROP COLUMN `client_ip`,
  DROP COLUMN `m_cpu_user_time_cnt`,
  DROP COLUMN `m_cpu_user_time_sum`,
  DROP COLUMN `m_cpu_sys_time_cnt`,
  DROP COLUMN `m_cpu_sys_time_sum`;

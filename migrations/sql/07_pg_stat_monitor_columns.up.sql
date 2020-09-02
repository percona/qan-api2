ALTER TABLE metrics
  ADD COLUMN `m_cpu_user_time_cnt` Float32,
  ADD COLUMN `m_cpu_user_time_sum` Float32,
  ADD COLUMN `m_cpu_sys_time_cnt` Float32,
  ADD COLUMN `m_cpu_sys_time_sum` Float32;

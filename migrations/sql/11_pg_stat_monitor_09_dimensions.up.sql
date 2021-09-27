ALTER TABLE metrics
  ADD COLUMN `top_queryid` HighCardinality(String),
  ADD COLUMN `application_name` NormalCardinality(String),
  ADD COLUMN `planid` HighCardinality(String);

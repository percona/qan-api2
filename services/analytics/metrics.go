// qan-api2
// Copyright (C) 2019 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package analitycs

import (
	"context"
	"fmt"
	"time"

	"github.com/percona/pmm/api/qanpb"
)

// GetMetrics implements rpc to get metrics for specific filtering.
func (s *Service) GetMetrics(ctx context.Context, in *qanpb.MetricsRequest) (*qanpb.MetricsReply, error) {
	fmt.Println("Call GetMetricsByDigest")
	labels := in.GetLabels()
	dbServers := []string{}
	dbSchemas := []string{}
	dbUsernames := []string{}
	clientHosts := []string{}
	dbLabels := map[string][]string{}
	for _, label := range labels {
		fmt.Printf("label: %v, : %v \n", label.Key, label.Value)
		switch label.Key {
		case "db_server":
			dbServers = label.Value
		case "db_schema":
			dbSchemas = label.Value
		case "db_username":
			dbUsernames = label.Value
		case "client_host":
			clientHosts = label.Value
		default:
			dbLabels[label.Key] = label.Value
		}
	}
	from := time.Unix(in.PeriodStartFrom.Seconds, 0)
	to := time.Unix(in.PeriodStartTo.Seconds, 0)
	m := make(map[string]*qanpb.MetricValues)
	resp := &qanpb.MetricsReply{
		Metrics: m,
	}
	metrics, err := s.mm.Get(
		ctx,
		from,
		to,
		in.FilterBy,
		in.GroupBy,
		dbServers,
		dbSchemas,
		dbUsernames,
		clientHosts,
		dbLabels,
	)
	if err != nil {
		return resp, fmt.Errorf("error in quering metrics:%v", err)
	}

	if len(metrics) > 2 {
		return resp, fmt.Errorf("not found for filter: %s and group: %s in given time range", in.FilterBy, in.GroupBy)
	}

	boolColumnNames := map[string]struct{}{
		"qc_hit":                 {},
		"full_scan":              {},
		"full_join":              {},
		"tmp_table":              {},
		"tmp_table_on_disk":      {},
		"filesort":               {},
		"filesort_on_disk":       {},
		"select_full_range_join": {},
		"select_range":           {},
		"select_range_check":     {},
		"sort_range":             {},
		"sort_rows":              {},
		"sort_scan":              {},
		"no_index_used":          {},
		"no_good_index_used":     {},
	}

	// "num_queries":           {},
	commonColumnNames := map[string]struct{}{
		"query_time":            {},
		"lock_time":             {},
		"rows_sent":             {},
		"rows_examined":         {},
		"rows_affected":         {},
		"rows_read":             {},
		"merge_passes":          {},
		"innodb_io_r_ops":       {},
		"innodb_io_r_bytes":     {},
		"innodb_io_r_wait":      {},
		"innodb_rec_lock_wait":  {},
		"innodb_queue_wait":     {},
		"innodb_pages_distinct": {},
		"query_length":          {},
		"bytes_sent":            {},
		"tmp_tables":            {},
		"tmp_disk_tables":       {},
		"tmp_table_sizes":       {},
		"docs_returned":         {},
		"response_length":       {},
		"docs_scanned":          {},
	}

	durationSec := to.Sub(from).Seconds()

	for k, _ := range commonColumnNames {
		cnt := interfaceToFloat32(metrics[0]["m_"+k+"_cnt"])
		sum := interfaceToFloat32(metrics[0]["m_"+k+"_sum"])
		totalSum := interfaceToFloat32(metrics[1]["m_"+k+"sum"])
		mv := qanpb.MetricValues{
			Cnt: cnt,
			Sum: sum,
			Min: interfaceToFloat32(metrics[0]["m_"+k+"_min"]),
			Max: interfaceToFloat32(metrics[0]["m_"+k+"_max"]),
			P99: interfaceToFloat32(metrics[0]["m_"+k+"_p99"]),
		}
		if cnt > 0 && sum > 0 {
			mv.Avg = sum / cnt
		}
		if sum > 0 && totalSum > 0 {
			mv.PTotal = sum / totalSum
		}
		if sum > 0 && durationSec > 0 {
			mv.Rate = sum / float32(durationSec)
		}
		resp.Metrics[k] = &mv
	}

	for k, _ := range boolColumnNames {
		cnt := interfaceToFloat32(metrics[0]["m_"+k+"_cnt"])
		sum := interfaceToFloat32(metrics[0]["m_"+k+"_sum"])
		totalSum := interfaceToFloat32(metrics[1]["m_"+k+"sum"])
		mv := qanpb.MetricValues{
			Cnt: cnt,
			Sum: sum,
		}
		if cnt > 0 && sum > 0 {
			mv.Avg = sum / cnt
		}
		if sum > 0 && totalSum > 0 {
			mv.PTotal = sum / totalSum
		}
		if sum > 0 && durationSec > 0 {
			mv.Rate = sum / float32(durationSec)
		}
		resp.Metrics[k] = &mv
	}

	return resp, err
}

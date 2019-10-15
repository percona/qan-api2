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

package models

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/percona/pmm/api/qanpb"
)

// Reporter implements models to select metrics bucket by params.
type Reporter struct {
	db *sqlx.DB
}

// NewReporter initialize Reporter with db instance.
func NewReporter(db *sqlx.DB) Reporter {
	return Reporter{db: db}
}

var funcMap = template.FuncMap{
	"inc":         func(i int) int { return i + 1 },
	"StringsJoin": strings.Join,
}

// M is map for interfaces.
type M map[string]interface{}

const queryReportTmpl = `
SELECT
{{ .Group }} AS dimension,
{{ if eq .Group "queryid" }} any(fingerprint) {{ else }} '' {{ end }} AS fingerprint,
{{range $j, $col := .CommonColumns}}
	SUM(m_{{ $col }}_cnt) AS m_{{ $col }}_cnt,
	SUM(m_{{ $col }}_sum) AS m_{{ $col }}_sum,
	MIN(m_{{ $col }}_min) AS m_{{ $col }}_min,
	MAX(m_{{ $col }}_max) AS m_{{ $col }}_max,
	AVG(m_{{ $col }}_p99) AS m_{{ $col }}_p99,
{{ end }}
{{range $j, $col := .SumColumns}}
        SUM(m_{{ $col }}_cnt) AS m_{{ $col }}_cnt,
        SUM(m_{{ $col }}_sum) AS m_{{ $col }}_sum,
{{ end }}
SUM(num_queries) AS num_queries,
{{range $j, $col := .SpecialColumns}}
    {{ if eq $col "load" }}
        {{ if $.IsQueryTimeInSelect }}
            m_query_time_sum / {{ $.PeriodDuration }} AS load,
        {{ else }}
            SUM(m_query_time_sum) / {{ $.PeriodDuration }} AS load,
        {{ end }}
	{{ else }}
		SUM({{ $col }}) AS {{ $col }},
	{{ end }}
{{ end }}
rowNumberInAllBlocks() AS total_rows
FROM metrics
WHERE period_start > :period_start_from AND period_start < :period_start_to
{{ if .Dimensions }}
    {{range $key, $vals := .Dimensions }}
        AND {{ $key }} IN ( '{{ StringsJoin $vals "', '" }}' )
    {{ end }}
{{ end }}
{{ if .Labels }}{{$i := 0}}
    AND ({{range $key, $vals := .Labels }}{{ $i = inc $i}}
        {{ if gt $i 1}} OR {{ end }} has(['{{ StringsJoin $vals "', '" }}'], labels.value[indexOf(labels.key, '{{ $key }}')])
    {{ end }})
{{ end }}
GROUP BY {{ .Group }}
        WITH TOTALS
ORDER BY {{ .Order }}
LIMIT :offset, :limit
`

//nolint
var tmplQueryReport = template.Must(template.New("queryReportTmpl").Funcs(funcMap).Parse(queryReportTmpl))

func inSlice(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// Select select metrics for report.
func (r *Reporter) Select(ctx context.Context, periodStartFromSec, periodStartToSec int64,
	dimensions map[string][]string, labels map[string][]string,
	group, order string, offset, limit uint32,
	specialColumns, commonColumns, sumColumns []string) ([]M, error) {

	arg := map[string]interface{}{
		"period_start_from": periodStartFromSec,
		"period_start_to":   periodStartToSec,
		"group":             group,
		"order":             order,
		"offset":            offset,
		"limit":             limit,
	}

	tmplArgs := struct {
		PeriodStartFrom     int64
		PeriodStartTo       int64
		PeriodDuration      int64
		Dimensions          map[string][]string
		Labels              map[string][]string
		Group               string
		Order               string
		Offset              uint32
		Limit               uint32
		SpecialColumns      []string
		CommonColumns       []string
		SumColumns          []string
		IsQueryTimeInSelect bool
	}{
		periodStartFromSec,
		periodStartToSec,
		periodStartToSec - periodStartFromSec,
		dimensions,
		labels,
		group,
		order,
		offset,
		limit,
		specialColumns,
		commonColumns,
		sumColumns,
		inSlice(commonColumns, "query_time"),
	}

	var queryBuffer bytes.Buffer

	if err := tmplQueryReport.Execute(&queryBuffer, tmplArgs); err != nil {
		return nil, errors.Wrap(err, "cannot execute tmplQueryReport")
	}

	var results []M
	query, args, err := sqlx.Named(queryBuffer.String(), arg)
	if err != nil {
		return nil, errors.Wrap(err, "prepare named tmplQueryReport")
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "populate arguments in IN clause")
	}
	query = r.db.Rebind(query)

	queryCtx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryxContext(queryCtx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "QueryxContext error")
	}
	for rows.Next() {
		result := make(M)
		err = rows.MapScan(result)
		if err != nil {
			return nil, errors.Wrap(err, "DimensionReport Scan error")
		}
		results = append(results, result)
	}
	rows.NextResultSet()
	total := make(M)
	for rows.Next() {
		err = rows.MapScan(total)
		if err != nil {
			return nil, errors.Wrap(err, "DimensionReport Scan TOTALS error")
		}
		results = append([]M{total}, results...)
	}
	return results, err
}

const queryReportSparklinesTmpl = `
SELECT
    intDivOrZero(toUnixTimestamp( :period_start_to ) - toUnixTimestamp(period_start), {{ .TimeFrame }}) AS point,
    toDateTime(toUnixTimestamp( :period_start_to ) - (point * {{ .TimeFrame }})) AS timestamp,
    {{ .TimeFrame }} AS time_frame,
    {{ if .IsCommon }}
        if(SUM(m_{{ .Column }}_cnt) == 0, NaN, SUM(m_{{ .Column }}_sum) / time_frame) AS m_{{ .Column }}_sum_per_sec
	{{ else }}
		{{ if eq .Column "num_queries" }}
			SUM(num_queries) / time_frame AS num_queries_per_sec
		{{ end }}
		{{ if eq .Column "num_queries_with_errors" }}
			SUM(num_queries_with_errors) / time_frame AS num_queries_with_errors_per_sec
		{{ end }}
		{{ if eq .Column "num_queries_with_warnings" }}
			SUM(num_queries_with_warnings) / time_frame AS num_queries_with_warnings_per_sec
		{{ end }}
		{{ if eq .Column "load" }}
			SUM(m_query_time_sum) / time_frame AS load
		{{ end }}
	{{ end }}
FROM metrics
WHERE period_start >= :period_start_from AND period_start <= :period_start_to
{{ if .DimensionVal }} AND {{ .Group }} = '{{ .DimensionVal }}' {{ end }}
    {{range $key, $vals := .Dimensions }} AND {{ $key }} IN ( '{{ StringsJoin $vals "', '" }}' ){{ end }}
{{ if .Labels }}{{$i := 0}}
    AND ({{range $key, $val := .Labels }} {{ $i = inc $i}}
        {{ if gt $i 1}} OR {{ end }} has(['{{ StringsJoin $val "', '" }}'], labels.value[indexOf(labels.key, '{{ $key }}')])
    {{ end }})
{{ end }}
GROUP BY point
ORDER BY point ASC;
`

//nolint
var tmplQueryReportSparklines = template.Must(template.New("queryReportSparklines").Funcs(funcMap).Parse(queryReportSparklinesTmpl))

// SelectSparklines selects datapoint for sparklines.
func (r *Reporter) SelectSparklines(ctx context.Context, dimensionVal string,
	periodStartFromSec, periodStartToSec int64,
	dimensions map[string][]string, labels map[string][]string,
	group string, column string) ([]*qanpb.Point, error) {

	// Align to minutes
	periodStartToSec = periodStartToSec / 60 * 60
	periodStartFromSec = periodStartFromSec / 60 * 60

	// If time range is bigger then two hour - amount of sparklines points = 120 to avoid huge data in response.
	// Otherwise amount of sparklines points is equal to minutes in in time range to not mess up calculation.
	amountOfPoints := int64(optimalAmountOfPoint)
	timePeriod := periodStartToSec - periodStartFromSec
	// reduce amount of point if period less then 2h.
	if timePeriod < int64((minFullTimeFrame).Seconds()) {
		// minimum point is 1 minute
		amountOfPoints = timePeriod / 60
	}

	// how many full minutes we can fit into given amount of points.
	minutesInPoint := (periodStartToSec - periodStartFromSec) / 60 / amountOfPoints
	// we need aditional point to show this minutes
	remainder := ((periodStartToSec - periodStartFromSec) / 60) % amountOfPoints
	amountOfPoints += remainder / minutesInPoint
	timeFrame := minutesInPoint * 60

	arg := map[string]interface{}{
		"dimension_val":     dimensionVal,
		"period_start_from": periodStartFromSec,
		"period_start_to":   periodStartToSec,
		"group":             group,
		"column":            column,
		"time_frame":        timeFrame,
	}

	tmplArgs := struct {
		DimensionVal    string
		PeriodStartFrom int64
		PeriodStartTo   int64
		PeriodDuration  int64
		Dimensions      map[string][]string
		Labels          map[string][]string
		Group           string
		Column          string
		IsCommon        bool
		TimeFrame       int64
	}{
		dimensionVal,
		periodStartFromSec,
		periodStartToSec,
		periodStartToSec - periodStartFromSec,
		dimensions,
		labels,
		group,
		column,
		!inSlice([]string{"load", "num_queries", "num_queries_with_errors", "num_queries_with_warnings"}, column),
		timeFrame,
	}

	var results []*qanpb.Point
	var queryBuffer bytes.Buffer

	if err := tmplQueryReportSparklines.Execute(&queryBuffer, tmplArgs); err != nil {
		return nil, errors.Wrap(err, "cannot execute tmplQueryReportSparklines")
	}
	query, args, err := sqlx.Named(queryBuffer.String(), arg)
	if err != nil {
		return nil, errors.Wrap(err, "prepare named")
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "populate arguments in IN clause")
	}
	query = r.db.Rebind(query)

	queryCtx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := r.db.QueryxContext(queryCtx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "report query")
	}
	resultsWithGaps := map[uint32]*qanpb.Point{}

	mainMetricColumnName := "m_query_time_sum"
	switch column {
	case "":
		mainMetricColumnName = "m_query_time_sum"
	case "load":
		mainMetricColumnName = "load"
	case "num_queries":
		mainMetricColumnName = "num_queries_per_sec"
	case "num_queries_with_errors":
		mainMetricColumnName = "num_queries_with_errors_per_sec"
	case "num_queries_with_warnings":
		mainMetricColumnName = "num_queries_with_warnings_per_sec"
	default:
		mainMetricColumnName = fmt.Sprintf("m_%s_sum_per_sec", column)
	}

	sparklinePointFieldsToQuery := []string{
		"point",
		"timestamp",
		"time_frame",
		mainMetricColumnName,
	}

	for rows.Next() {
		p := qanpb.Point{}
		res := getPointFieldsList(&p, sparklinePointFieldsToQuery)
		err = rows.Scan(res...)
		if err != nil {
			return nil, errors.Wrap(err, "SelectSparklines scan errors")
		}
		resultsWithGaps[p.Point] = &p
	}

	// fill in gaps in time series.
	for pointN := uint32(0); int64(pointN) < amountOfPoints; pointN++ {
		p, ok := resultsWithGaps[pointN]
		if !ok {
			p = &qanpb.Point{}
			p.Point = pointN
			p.TimeFrame = uint32(timeFrame)
			timeShift := timeFrame * int64(pointN)
			ts := periodStartToSec - timeShift
			p.Timestamp = time.Unix(ts, 0).UTC().Format(time.RFC3339)
		}
		results = append(results, p)
	}

	return results, err
}

const queryDimension = `
SELECT '{{ .DimensionName }}' AS key, {{ .DimensionName }} AS value, SUM( {{ .MainMetric }} ) AS main_metric_sum
  FROM metrics
 WHERE period_start >= ?
   AND period_start <= ?
   {{range $key, $vals := .Dimensions }} AND {{ $key }} IN ( '{{ StringsJoin $vals "', '" }}' ){{ end }}
    {{ if .Labels }}{{$i := 0}}
       AND ({{range $key, $val := .Labels }} {{ $i = inc $i}}
                {{ if gt $i 1}} AND {{ end }} has(['{{ StringsJoin $val "', '" }}'], labels.value[indexOf(labels.key, '{{ $key }}')])
           {{ end }})
    {{ end }}
GROUP BY {{ .DimensionName }}
 WITH TOTALS
ORDER BY main_metric_sum DESC, value;
`

const queryLabels = `
SELECT arrayJoin(labels.key) AS key, labels.value[indexOf(labels.key, key)] AS value, SUM( {{ .MainMetric }} ) AS main_metric_sum
  FROM metrics
 WHERE period_start >= ?
   AND period_start <= ?
    {{range $keyD, $valsD := .Dimensions }} AND {{ $keyD }} IN ( '{{ StringsJoin $valsD "', '" }}' ){{ end }}
    {{ if .Labels }}{{$i := 0}} AND
        ({{range $keyL, $valsL := .Labels }} {{ $i = inc $i}}
            {{ if gt $i 1}} AND {{ end }} has(['{{ StringsJoin $valsL "', '" }}'], labels.value[indexOf(labels.key, '{{ $keyL }}')])
        {{ end }})
    {{ end }}
GROUP BY key, value
ORDER BY main_metric_sum DESC, key, value;
`

type customLabel struct {
	key              string
	value            string
	mainMetricPerSec float32
}

var queryDimensionTmpl = template.Must(template.New("queryDimension").Funcs(funcMap).Parse(queryDimension))
var queryLabelsTmpl = template.Must(template.New("queryLabels").Funcs(funcMap).Parse(queryLabels))
var dimensionQueries = map[string]*template.Template{
	"service_name":    queryDimensionTmpl,
	"database":        queryDimensionTmpl,
	"schema":          queryDimensionTmpl,
	"username":        queryDimensionTmpl,
	"client_host":     queryDimensionTmpl,
	"replication_set": queryDimensionTmpl,
	"cluster":         queryDimensionTmpl,
	"service_type":    queryDimensionTmpl,
	"service_id":      queryDimensionTmpl,
	"environment":     queryDimensionTmpl,
	"az":              queryDimensionTmpl,
	"region":          queryDimensionTmpl,
	"node_model":      queryDimensionTmpl,
	"node_id":         queryDimensionTmpl,
	"node_name":       queryDimensionTmpl,
	"node_type":       queryDimensionTmpl,
	"machine_id":      queryDimensionTmpl,
	"container_name":  queryDimensionTmpl,
	"container_id":    queryDimensionTmpl,
	"labels":          queryLabelsTmpl,
}

// SelectFilters selects dimension and their values, and also keys and values of labels.
func (r *Reporter) SelectFilters(ctx context.Context, periodStartFromSec, periodStartToSec int64, mainMetricName string, dimensions, labels map[string][]string) (*qanpb.FiltersReply, error) {
	result := qanpb.FiltersReply{
		Labels: make(map[string]*qanpb.ListLabels),
	}

	if !isValidMetricColumn(mainMetricName) {
		return nil, fmt.Errorf("invalid main metric name %s", mainMetricName)
	}

	for dimensionName, dimensionQuery := range dimensionQueries {
		values, mainMetricPerSec, err := r.queryFilters(ctx, periodStartFromSec, periodStartToSec, dimensionName, mainMetricName, dimensionQuery, dimensions, labels)
		if err != nil {
			return nil, errors.Wrap(err, "cannot select "+dimensionName+" dimension")
		}

		totals := map[string]float32{}
		if mainMetricPerSec == 0 {
			for _, label := range values {
				totals[label.key] += label.mainMetricPerSec
			}
		}

		for _, label := range values {
			if _, ok := result.Labels[label.key]; !ok {
				result.Labels[label.key] = &qanpb.ListLabels{
					Name: []*qanpb.Values{},
				}
			}
			total := mainMetricPerSec
			if mainMetricPerSec == 0 {
				total = totals[label.key]
			}
			val := qanpb.Values{
				Value:             label.value,
				MainMetricPerSec:  label.mainMetricPerSec,
				MainMetricPercent: label.mainMetricPerSec / total,
			}
			result.Labels[label.key].Name = append(result.Labels[label.key].Name, &val)
		}
	}

	return &result, nil
}

func (r *Reporter) queryFilters(ctx context.Context, periodStartFromSec,
	periodStartToSec int64, dimensionName, mainMetricName string, tmplQueryFilter *template.Template, queryDimensions, queryLabels map[string][]string) ([]*customLabel, float32, error) {
	durationSec := periodStartToSec - periodStartFromSec
	var labels []*customLabel

	tmplArgs := struct {
		MainMetric    string
		DimensionName string
		Dimensions    map[string][]string
		Labels        map[string][]string
	}{
		mainMetricName,
		dimensionName,
		queryDimensions,
		queryLabels,
	}

	var queryBuffer bytes.Buffer

	if err := tmplQueryFilter.Execute(&queryBuffer, tmplArgs); err != nil {
		return nil, 0, errors.Wrap(err, "cannot execute tmplQueryFilter"+queryBuffer.String())
	}

	rows, err := r.db.QueryContext(ctx, queryBuffer.String(), periodStartFromSec, periodStartToSec)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to select for QueryFilter "+queryBuffer.String())
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var label customLabel
		err = rows.Scan(&label.key, &label.value, &label.mainMetricPerSec)
		if err != nil {
			return nil, 0, errors.Wrap(err, "failed to scan for QueryFilter "+queryBuffer.String())
		}
		label.mainMetricPerSec /= float32(durationSec)
		labels = append(labels, &label)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "failed to select for QueryFilter "+queryBuffer.String())
	}

	totalMainMetricPerSec := float32(0)

	if rows.NextResultSet() {
		var labelTotal customLabel
		for rows.Next() {
			err = rows.Scan(&labelTotal.key, &labelTotal.value, &labelTotal.mainMetricPerSec)
			if err != nil {
				return nil, 0, errors.Wrap(err, "failed to scan total for QueryFilter "+queryBuffer.String())
			}
			totalMainMetricPerSec = labelTotal.mainMetricPerSec / float32(durationSec)
		}
		if err = rows.Err(); err != nil {
			return nil, 0, errors.Wrap(err, "failed to select total for QueryFilter "+queryBuffer.String())
		}
	}

	return labels, totalMainMetricPerSec, nil
}

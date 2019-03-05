// qan-api
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

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/percona/go-mysql/event"
	slowlog "github.com/percona/go-mysql/log"
	parser "github.com/percona/go-mysql/log/slow"
	"github.com/percona/go-mysql/query"
	pbqan "github.com/percona/pmm/api/qan"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const agentUUID = "dc889ca7be92a66f0a00f616f69ffa7b"
const dbServer = "fb_db"

type closedChannelError struct {
	error
}

type QueryClassDimentions struct {
	DbUsername  string
	ClientHost  string
	PeriodStart int64
	PeriodEnd   int64
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	logOpt := slowlog.Options{}
	slowLogPath := flag.String("slow-log", "logs/mysql-slow.log", "Path to MySQL slow log file")
	logTimeStart := flag.String("logTimeStart", "2019-01-01 00:00:00", "Start fake time of query from")
	serverURL := flag.String("server-url", "127.0.0.1:80", "ULR of QAN-API Server")
	newEventWait := flag.Duration("new-event-wait", 10*time.Second, "Time to wait for a new event in slow log.")
	flag.Parse()

	log.SetOutput(os.Stderr)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*serverURL, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	client := pbqan.NewAgentClient(conn)

	events := parseSlowLog(*slowLogPath, logOpt)

	ctx := context.TODO()
	stream, err := client.DataInterchange(ctx)
	if err != nil {
		log.Fatalf("%v.DataInterchange(), %v", client, err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				continue
			}
			log.Printf("Got message from Api. Saved %d query class(es)", in.SavedAmount)
		}
	}()

	events = parseSlowLog(*slowLogPath, logOpt)
	fmt.Println("Parsing slowlog: ", *slowLogPath, "...")

	logStart, _ := time.Parse("2006-01-02 15:04:05", *logTimeStart)
	periodNumber := 0

	var periodStart time.Time
	if periodStart.IsZero() {
		periodStart = logStart.Add(time.Duration(periodNumber) * time.Minute)
	}

	for {
		// start := time.Now()
		var prewTs time.Time
		err = bulkSend(stream, func(am *pbqan.AgentMessage) error {
			i := 0
			aggregator := event.NewAggregator(true, 0, 1) // add right params
			for e := range events {
				fingerprint := query.Fingerprint(e.Query)
				digest := query.Id(fingerprint)

				e.Db = fmt.Sprintf("schema%d", r.Intn(100))      // fake data 100
				e.User = fmt.Sprintf("user%d", r.Intn(100))      // fake data 100
				e.Host = fmt.Sprintf("10.11.12.%d", r.Intn(100)) // fake data 100
				e.Server = fmt.Sprintf("db%d", r.Intn(10))       // fake data 10
				e.LabelsKey = []string{fmt.Sprintf("label%d", r.Intn(10)), fmt.Sprintf("label%d", r.Intn(10)), fmt.Sprintf("label%d", r.Intn(10))}
				e.LabelsValue = []string{fmt.Sprintf("value%d", r.Intn(100)), fmt.Sprintf("value%d", r.Intn(100)), fmt.Sprintf("value%d", r.Intn(100))}

				aggregator.AddEvent(e, digest, e.User, e.Host, e.Db, e.Server, fingerprint)
				i++

				// Pass last offset to restart reader when reached out end of slowlog.
				logOpt.StartOffset = e.OffsetEnd

				if prewTs.IsZero() {
					prewTs = e.Ts
				}

				if e.Ts.Sub(prewTs).Seconds() > 59 {
					prewTs = e.Ts
					periodStart = logStart.Add(time.Duration(periodNumber) * time.Minute)
					periodNumber++
					break
				}
			}

			res := aggregator.Finalize()

			for _, v := range res.Class {

				qc := &pbqan.QueryClass{
					Queryid:       v.Id,
					Fingerprint:   v.Fingerprint,
					DDatabase:     "",
					DSchema:       v.Db,
					DUsername:     v.User,
					DClientHost:   v.Host,
					DServer:       v.Server,
					Labels:        listsToMap(v.LabelsKey, v.LabelsValue),
					AgentUuid:     agentUUID,
					PeriodStart:   periodStart.Truncate(1 * time.Minute).Unix(),
					PeriodLength:  uint32(60),
					Example:       v.Example.Query,
					ExampleFormat: 1,
					ExampleType:   1,
					NumQueries:    float32(v.TotalQueries),
				}

				// If key has suffix _time or _wait than field is TimeMetrics.
				// For Boolean metrics exists only Sum.
				// https://www.percona.com/doc/percona-server/5.7/diagnostics/slow_extended.html
				// TimeMetrics: query_time, lock_time, rows_sent, innodb_io_r_wait, innodb_rec_lock_wait, innodb_queue_wait.
				// NumberMetrics: rows_examined, rows_affected, rows_read, merge_passes, innodb_io_r_ops, innodb_io_r_bytes,
				// innodb_pages_distinct, query_length, bytes_sent, tmp_tables, tmp_disk_tables, tmp_table_sizes.
				// BooleanMetrics: qc_hit, full_scan, full_join, tmp_table, tmp_table_on_disk, filesort, filesort_on_disk,
				// select_full_range_join, select_range, select_range_check, sort_range, sort_rows, sort_scan,
				// no_index_used, no_good_index_used.

				// query_time - Query_time
				if m, ok := v.Metrics.TimeMetrics["Query_time"]; ok {
					qc.MQueryTimeSum = float32(m.Sum)
					qc.MQueryTimeMax = float32(*m.Max)
					qc.MQueryTimeMin = float32(*m.Min)
					qc.MQueryTimeP99 = float32(*m.P99)
				}
				// lock_time - Lock_time
				if m, ok := v.Metrics.TimeMetrics["Lock_time"]; ok {
					qc.MLockTimeSum = float32(m.Sum)
					qc.MLockTimeMax = float32(*m.Max)
					qc.MLockTimeMin = float32(*m.Min)
					qc.MLockTimeP99 = float32(*m.P99)
				}
				// rows_sent - Rows_sent
				if m, ok := v.Metrics.NumberMetrics["Rows_sent"]; ok {
					qc.MRowsSentSum = float32(m.Sum)
					qc.MRowsSentMax = float32(*m.Max)
					qc.MRowsSentMin = float32(*m.Min)
					qc.MRowsSentP99 = float32(*m.P99)
				}
				// rows_examined - Rows_examined
				if m, ok := v.Metrics.NumberMetrics["Rows_examined"]; ok {
					qc.MRowsExaminedSum = float32(m.Sum)
					qc.MRowsExaminedMax = float32(*m.Max)
					qc.MRowsExaminedMin = float32(*m.Min)
					qc.MRowsExaminedP99 = float32(*m.P99)
				}
				// rows_affected - Rows_affected
				if m, ok := v.Metrics.NumberMetrics["Rows_affected"]; ok {
					qc.MRowsAffectedSum = float32(m.Sum)
					qc.MRowsAffectedMax = float32(*m.Max)
					qc.MRowsAffectedMin = float32(*m.Min)
					qc.MRowsAffectedP99 = float32(*m.P99)
				}
				// rows_read - Rows_read
				if m, ok := v.Metrics.NumberMetrics["Rows_read"]; ok {
					qc.MRowsReadSum = float32(m.Sum)
					qc.MRowsReadMax = float32(*m.Max)
					qc.MRowsReadMin = float32(*m.Min)
					qc.MRowsReadP99 = float32(*m.P99)
				}
				// merge_passes - Merge_passes
				if m, ok := v.Metrics.NumberMetrics["Merge_passes"]; ok {
					qc.MMergePassesSum = float32(m.Sum)
					qc.MMergePassesMax = float32(*m.Max)
					qc.MMergePassesMin = float32(*m.Min)
					qc.MMergePassesP99 = float32(*m.P99)
				}
				// innodb_io_r_ops - InnoDB_IO_r_ops
				if m, ok := v.Metrics.NumberMetrics["InnoDB_IO_r_ops"]; ok {
					qc.MInnodbIoROpsSum = float32(m.Sum)
					qc.MInnodbIoROpsMax = float32(*m.Max)
					qc.MInnodbIoROpsMin = float32(*m.Min)
					qc.MInnodbIoROpsP99 = float32(*m.P99)
				}
				// innodb_io_r_bytes - InnoDB_IO_r_bytes
				if m, ok := v.Metrics.NumberMetrics["InnoDB_IO_r_bytes"]; ok {
					qc.MInnodbIoRBytesSum = float32(m.Sum)
					qc.MInnodbIoRBytesMax = float32(*m.Max)
					qc.MInnodbIoRBytesMin = float32(*m.Min)
					qc.MInnodbIoRBytesP99 = float32(*m.P99)
				}
				// innodb_io_r_wait - InnoDB_IO_r_wait
				if m, ok := v.Metrics.TimeMetrics["InnoDB_IO_r_wait"]; ok {
					qc.MInnodbIoRWaitSum = float32(m.Sum)
					qc.MInnodbIoRWaitMax = float32(*m.Max)
					qc.MInnodbIoRWaitMin = float32(*m.Min)
					qc.MInnodbIoRWaitP99 = float32(*m.P99)
				}
				// innodb_rec_lock_wait - InnoDB_rec_lock_wait
				if m, ok := v.Metrics.TimeMetrics["InnoDB_rec_lock_wait"]; ok {
					qc.MInnodbRecLockWaitSum = float32(m.Sum)
					qc.MInnodbRecLockWaitMax = float32(*m.Max)
					qc.MInnodbRecLockWaitMin = float32(*m.Min)
					qc.MInnodbRecLockWaitP99 = float32(*m.P99)
				}
				// innodb_queue_wait - InnoDB_queue_wait
				if m, ok := v.Metrics.TimeMetrics["InnoDB_queue_wait"]; ok {
					qc.MInnodbQueueWaitSum = float32(m.Sum)
					qc.MInnodbQueueWaitMax = float32(*m.Max)
					qc.MInnodbQueueWaitMin = float32(*m.Min)
					qc.MInnodbQueueWaitP99 = float32(*m.P99)
				}
				// innodb_pages_distinct - InnoDB_pages_distinct
				if m, ok := v.Metrics.NumberMetrics["InnoDB_pages_distinct"]; ok {
					qc.MInnodbPagesDistinctSum = float32(m.Sum)
					qc.MInnodbPagesDistinctMax = float32(*m.Max)
					qc.MInnodbPagesDistinctMin = float32(*m.Min)
					qc.MInnodbPagesDistinctP99 = float32(*m.P99)
				}
				// query_length - Query_length
				if m, ok := v.Metrics.NumberMetrics["Query_length"]; ok {
					qc.MQueryLengthSum = float32(m.Sum)
					qc.MQueryLengthMax = float32(*m.Max)
					qc.MQueryLengthMin = float32(*m.Min)
					qc.MQueryLengthP99 = float32(*m.P99)
				}
				// bytes_sent - Bytes_sent
				if m, ok := v.Metrics.NumberMetrics["Bytes_sent"]; ok {
					qc.MBytesSentSum = float32(m.Sum)
					qc.MBytesSentMax = float32(*m.Max)
					qc.MBytesSentMin = float32(*m.Min)
					qc.MBytesSentP99 = float32(*m.P99)
				}
				// tmp_tables - Tmp_tables
				if m, ok := v.Metrics.NumberMetrics["Tmp_tables"]; ok {
					qc.MTmpTablesSum = float32(m.Sum)
					qc.MTmpTablesMax = float32(*m.Max)
					qc.MTmpTablesMin = float32(*m.Min)
					qc.MTmpTablesP99 = float32(*m.P99)
				}
				// tmp_disk_tables - Tmp_disk_tables
				if m, ok := v.Metrics.NumberMetrics["Tmp_disk_tables"]; ok {
					qc.MTmpDiskTablesSum = float32(m.Sum)
					qc.MTmpDiskTablesMax = float32(*m.Max)
					qc.MTmpDiskTablesMin = float32(*m.Min)
					qc.MTmpDiskTablesP99 = float32(*m.P99)
				}
				// tmp_table_sizes - Tmp_table_sizes
				if m, ok := v.Metrics.NumberMetrics["Tmp_table_sizes"]; ok {
					qc.MTmpTableSizesSum = float32(m.Sum)
					qc.MTmpTableSizesMax = float32(*m.Max)
					qc.MTmpTableSizesMin = float32(*m.Min)
					qc.MTmpTableSizesP99 = float32(*m.P99)
				}
				// qc_hit - QC_Hit
				if m, ok := v.Metrics.BoolMetrics["QC_Hit"]; ok {
					qc.MQcHitCnt = float32(m.Cnt)
					qc.MQcHitSum = float32(m.Sum)
				}
				// full_scan - Full_scan
				if m, ok := v.Metrics.BoolMetrics["Full_scan"]; ok {
					qc.MFullScanCnt = float32(m.Cnt)
					qc.MFullScanSum = float32(m.Sum)
				}
				// full_join - Full_join
				if m, ok := v.Metrics.BoolMetrics["Full_join"]; ok {
					qc.MFullJoinCnt = float32(m.Cnt)
					qc.MFullJoinSum = float32(m.Sum)
				}
				// tmp_table - Tmp_table
				if m, ok := v.Metrics.BoolMetrics["Tmp_table"]; ok {
					qc.MTmpTableCnt = float32(m.Cnt)
					qc.MTmpTableSum = float32(m.Sum)
				}
				// tmp_table_on_disk - Tmp_table_on_disk
				if m, ok := v.Metrics.BoolMetrics["Tmp_table_on_disk"]; ok {
					qc.MTmpTableOnDiskCnt = float32(m.Cnt)
					qc.MTmpTableOnDiskSum = float32(m.Sum)
				}
				// filesort - Filesort
				if m, ok := v.Metrics.BoolMetrics["Filesort"]; ok {
					qc.MFilesortCnt = float32(m.Cnt)
					qc.MFilesortSum = float32(m.Sum)
				}
				// filesort_on_disk - Filesort_on_disk
				if m, ok := v.Metrics.BoolMetrics["Filesort_on_disk"]; ok {
					qc.MFilesortOnDiskCnt = float32(m.Cnt)
					qc.MFilesortOnDiskSum = float32(m.Sum)
				}
				// select_full_range_join - Select_full_range_join
				if m, ok := v.Metrics.BoolMetrics["Select_full_range_join"]; ok {
					qc.MSelectFullRangeJoinCnt = float32(m.Cnt)
					qc.MSelectFullRangeJoinSum = float32(m.Sum)
				}
				// select_range - Select_range
				if m, ok := v.Metrics.BoolMetrics["Select_range"]; ok {
					qc.MSelectRangeCnt = float32(m.Cnt)
					qc.MSelectRangeSum = float32(m.Sum)
				}
				// select_range_check - Select_range_check
				if m, ok := v.Metrics.BoolMetrics["Select_range_check"]; ok {
					qc.MSelectRangeCheckCnt = float32(m.Cnt)
					qc.MSelectRangeCheckSum = float32(m.Sum)
				}
				// sort_range - Sort_range
				if m, ok := v.Metrics.BoolMetrics["Sort_range"]; ok {
					qc.MSortRangeCnt = float32(m.Cnt)
					qc.MSortRangeSum = float32(m.Sum)
				}
				// sort_rows - Sort_rows
				if m, ok := v.Metrics.BoolMetrics["Sort_rows"]; ok {
					qc.MSortRowsCnt = float32(m.Cnt)
					qc.MSortRowsSum = float32(m.Sum)
				}
				// sort_scan - Sort_scan
				if m, ok := v.Metrics.BoolMetrics["Sort_scan"]; ok {
					qc.MSortScanCnt = float32(m.Cnt)
					qc.MSortScanSum = float32(m.Sum)
				}
				// no_index_used - No_index_used
				if m, ok := v.Metrics.BoolMetrics["No_index_used"]; ok {
					qc.MNoIndexUsedCnt = float32(m.Cnt)
					qc.MNoIndexUsedSum = float32(m.Sum)
				}
				// no_good_index_used - No_good_index_used
				if m, ok := v.Metrics.BoolMetrics["No_good_index_used"]; ok {
					qc.MNoGoodIndexUsedCnt = float32(m.Cnt)
					qc.MNoGoodIndexUsedSum = float32(m.Sum)
				}

				am.QueryClass = append(am.QueryClass, qc)
			}

			// No new events in slowlog. Nothing to send to API.
			if i == 0 {
				return closedChannelError{errors.New("closed channel")}
			}
			// Reached end of slowlog. Send all what we have to API.
			return nil
		})

		if err != nil {
			if _, ok := err.(closedChannelError); !ok {
				log.Fatal("transaction error:", err)
			}
			// Channel is closed when reached end of the slowlog.
			// Wait and try read the slowlog again.
			time.Sleep(*newEventWait)
			events = parseSlowLog(*slowLogPath, logOpt)
		}
	}
}

func bulkSend(stream pbqan.Agent_DataInterchangeClient, fn func(*pbqan.AgentMessage) error) error {
	am := &pbqan.AgentMessage{}
	err := fn(am)
	if err != nil {
		return err
	}
	lenQC := len(am.QueryClass)
	if lenQC > 0 {
		if err := stream.Send(am); err != nil {
			return fmt.Errorf("sent error: %v", err)
		}
		fmt.Printf("send to qan %v QC\n", lenQC)
	}
	return nil
}

func parseSlowLog(filename string, o slowlog.Options) <-chan *slowlog.Event {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		log.Fatal("cannot open slowlog. Use --slow-log=/path/to/slow.log", err)
	}
	p := parser.NewSlowLogParser(file, o)
	go func() {
		err = p.Start()
		if err != nil {
			log.Fatal("cannot start parser", err)
		}
	}()
	return p.EventChan()
}

func listsToMap(k, v []string) map[string]string {
	m := map[string]string{}
	for i, e := range k {
		m[e] = v[i]
	}
	return m
}

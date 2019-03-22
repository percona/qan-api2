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
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	_ "github.com/kshvakov/clickhouse"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jmoiron/sqlx"
	"github.com/percona/pmm/api/qanpb"
	"github.com/percona/qan-api2/models"
)

func TestService_GetFilters(t *testing.T) {

	dsn, ok := os.LookupEnv("QANAPI_DSN_TEST")
	if !ok {
		dsn = "clickhouse://127.0.0.1:19000?database=pmm_test"
	}
	db, err := sqlx.Connect("clickhouse", dsn)
	if err != nil {
		log.Fatal("Connection: ", err)
	}

	rm := models.NewReporter(db)
	mm := models.NewMetrics(db)
	t1, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2019-01-01T00:01:00Z")
	var want qanpb.FiltersReply
	err = json.Unmarshal([]byte(expectedJSON), &want)

	type fields struct {
		rm models.Reporter
		mm models.Metrics
	}
	type args struct {
		ctx context.Context
		in  *qanpb.FiltersRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *qanpb.FiltersReply
		wantErr bool
	}{
		{
			"success",
			fields{rm: rm, mm: mm},
			args{
				context.TODO(),
				&qanpb.FiltersRequest{
					From: &timestamp.Timestamp{Seconds: t1.Unix()},
					To:   &timestamp.Timestamp{Seconds: t2.Unix()},
				},
			},
			&want,
			false,
		},
		{
			"fail",
			fields{rm: rm, mm: mm},
			args{
				context.TODO(),
				&qanpb.FiltersRequest{
					From: &timestamp.Timestamp{Seconds: t2.Unix()},
					To:   &timestamp.Timestamp{Seconds: t1.Unix()},
				},
			},
			&qanpb.FiltersReply{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				rm: tt.fields.rm,
				mm: tt.fields.mm,
			}
			got, err := s.GetFilters(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

const expectedJSON = `
{
	"labels": {
		"d_client_host": {
			"values": [
				{
					"value": "10.11.12.96",
					"count": 1
				},
				{
					"value": "10.11.12.65",
					"count": 1
				},
				{
					"value": "10.11.12.43",
					"count": 1
				},
				{
					"value": "10.11.12.49",
					"count": 1
				},
				{
					"value": "10.11.12.23",
					"count": 1
				},
				{
					"value": "10.11.12.32",
					"count": 1
				},
				{
					"value": "10.11.12.33",
					"count": 2
				},
				{
					"value": "10.11.12.15",
					"count": 1
				},
				{
					"value": "10.11.12.14",
					"count": 1
				},
				{
					"value": "10.11.12.10",
					"count": 1
				},
				{
					"value": "10.11.12.11",
					"count": 2
				},
				{
					"value": "10.11.12.89",
					"count": 1
				},
				{
					"value": "10.11.12.36",
					"count": 2
				},
				{
					"value": "10.11.12.71",
					"count": 1
				},
				{
					"value": "10.11.12.83",
					"count": 2
				},
				{
					"value": "10.11.12.8",
					"count": 3
				},
				{
					"value": "10.11.12.87",
					"count": 1
				},
				{
					"value": "10.11.12.53",
					"count": 1
				},
				{
					"value": "10.11.12.52",
					"count": 1
				},
				{
					"value": "10.11.12.74",
					"count": 2
				},
				{
					"value": "10.11.12.85",
					"count": 1
				},
				{
					"value": "10.11.12.55",
					"count": 1
				},
				{
					"value": "10.11.12.54",
					"count": 1
				},
				{
					"value": "10.11.12.73",
					"count": 2
				},
				{
					"value": "10.11.12.81",
					"count": 1
				},
				{
					"value": "10.11.12.18",
					"count": 1
				},
				{
					"value": "10.11.12.1",
					"count": 1
				},
				{
					"value": "10.11.12.13",
					"count": 1
				},
				{
					"value": "10.11.12.12",
					"count": 1
				},
				{
					"value": "10.11.12.79",
					"count": 1
				},
				{
					"value": "10.11.12.16",
					"count": 1
				},
				{
					"value": "10.11.12.4",
					"count": 1
				},
				{
					"value": "10.11.12.20",
					"count": 1
				},
				{
					"value": "10.11.12.25",
					"count": 1
				},
				{
					"value": "10.11.12.24",
					"count": 1
				},
				{
					"value": "10.11.12.69",
					"count": 1
				},
				{
					"value": "10.11.12.68",
					"count": 1
				},
				{
					"value": "10.11.12.44",
					"count": 1
				},
				{
					"value": "10.11.12.45",
					"count": 1
				},
				{
					"value": "10.11.12.95",
					"count": 2
				}
			]
		},
		"d_database": {
			"values": [
				{
					"value": "schema52",
					"count": 1
				},
				{
					"value": "schema21",
					"count": 1
				},
				{
					"value": "schema95",
					"count": 1
				},
				{
					"value": "schema73",
					"count": 1
				},
				{
					"value": "schema25",
					"count": 1
				},
				{
					"value": "schema36",
					"count": 1
				},
				{
					"value": "schema98",
					"count": 1
				},
				{
					"value": "schema45",
					"count": 2
				},
				{
					"value": "schema17",
					"count": 1
				},
				{
					"value": "schema82",
					"count": 1
				},
				{
					"value": "schema91",
					"count": 1
				},
				{
					"value": "schema20",
					"count": 2
				},
				{
					"value": "schema53",
					"count": 1
				},
				{
					"value": "schema33",
					"count": 1
				},
				{
					"value": "schema99",
					"count": 1
				},
				{
					"value": "schema37",
					"count": 2
				},
				{
					"value": "schema76",
					"count": 2
				},
				{
					"value": "schema90",
					"count": 1
				},
				{
					"value": "schema65",
					"count": 1
				},
				{
					"value": "schema71",
					"count": 2
				},
				{
					"value": "schema3",
					"count": 1
				},
				{
					"value": "schema39",
					"count": 1
				},
				{
					"value": "schema8",
					"count": 1
				},
				{
					"value": "schema59",
					"count": 1
				},
				{
					"value": "schema84",
					"count": 1
				},
				{
					"value": "schema62",
					"count": 1
				},
				{
					"value": "schema30",
					"count": 1
				},
				{
					"value": "schema1",
					"count": 2
				},
				{
					"value": "schema78",
					"count": 1
				},
				{
					"value": "schema18",
					"count": 1
				},
				{
					"value": "schema75",
					"count": 1
				},
				{
					"value": "schema93",
					"count": 2
				},
				{
					"value": "schema9",
					"count": 1
				},
				{
					"value": "schema10",
					"count": 1
				},
				{
					"value": "schema38",
					"count": 2
				},
				{
					"value": "schema70",
					"count": 1
				},
				{
					"value": "schema6",
					"count": 1
				},
				{
					"value": "schema79",
					"count": 1
				},
				{
					"value": "schema42",
					"count": 1
				},
				{
					"value": "schema74",
					"count": 1
				},
				{
					"value": "schema35",
					"count": 1
				}
			]
		},
		"d_schema": {
			"values": [
				{
					"count": 49
				}
			]
		},
		"d_server": {
			"values": [
				{
					"value": "db3",
					"count": 7
				},
				{
					"value": "db1",
					"count": 4
				},
				{
					"value": "db7",
					"count": 2
				},
				{
					"value": "db6",
					"count": 8
				},
				{
					"value": "db4",
					"count": 3
				},
				{
					"value": "db8",
					"count": 7
				},
				{
					"value": "db0",
					"count": 5
				},
				{
					"value": "db2",
					"count": 3
				},
				{
					"value": "db5",
					"count": 6
				},
				{
					"value": "db9",
					"count": 4
				}
			]
		},
		"d_username": {
			"values": [
				{
					"value": "user2",
					"count": 2
				},
				{
					"value": "user16",
					"count": 1
				},
				{
					"value": "user90",
					"count": 1
				},
				{
					"value": "user78",
					"count": 1
				},
				{
					"value": "user70",
					"count": 1
				},
				{
					"value": "user44",
					"count": 2
				},
				{
					"value": "user54",
					"count": 1
				},
				{
					"value": "user63",
					"count": 2
				},
				{
					"value": "user98",
					"count": 2
				},
				{
					"value": "user82",
					"count": 1
				},
				{
					"value": "user1",
					"count": 1
				},
				{
					"value": "user43",
					"count": 2
				},
				{
					"value": "user64",
					"count": 1
				},
				{
					"value": "user20",
					"count": 1
				},
				{
					"value": "user38",
					"count": 3
				},
				{
					"value": "user52",
					"count": 1
				},
				{
					"value": "user47",
					"count": 1
				},
				{
					"value": "user69",
					"count": 2
				},
				{
					"value": "user22",
					"count": 1
				},
				{
					"value": "user71",
					"count": 1
				},
				{
					"value": "user55",
					"count": 1
				},
				{
					"value": "user96",
					"count": 1
				},
				{
					"value": "user5",
					"count": 1
				},
				{
					"value": "user76",
					"count": 1
				},
				{
					"value": "user17",
					"count": 1
				},
				{
					"value": "user45",
					"count": 1
				},
				{
					"value": "user36",
					"count": 1
				},
				{
					"value": "user35",
					"count": 1
				},
				{
					"value": "user74",
					"count": 1
				},
				{
					"value": "user67",
					"count": 1
				},
				{
					"value": "user59",
					"count": 2
				},
				{
					"value": "user31",
					"count": 2
				},
				{
					"value": "user84",
					"count": 1
				},
				{
					"value": "user62",
					"count": 1
				},
				{
					"value": "user92",
					"count": 1
				},
				{
					"value": "user97",
					"count": 1
				},
				{
					"value": "user9",
					"count": 1
				},
				{
					"value": "user21",
					"count": 1
				},
				{
					"value": "user32",
					"count": 1
				}
			]
		},
		"label0": {
			"values": [
				{
					"value": "value1",
					"count": 1
				},
				{
					"value": "value39",
					"count": 1
				},
				{
					"value": "value43",
					"count": 1
				},
				{
					"value": "value51",
					"count": 1
				},
				{
					"value": "value55",
					"count": 1
				},
				{
					"value": "value57",
					"count": 1
				},
				{
					"value": "value62",
					"count": 1
				},
				{
					"value": "value68",
					"count": 1
				},
				{
					"value": "value82",
					"count": 1
				},
				{
					"value": "value88",
					"count": 1
				},
				{
					"value": "value92",
					"count": 1
				},
				{
					"value": "value95",
					"count": 1
				}
			]
		},
		"label1": {
			"values": [
				{
					"value": "value0",
					"count": 1
				},
				{
					"value": "value10",
					"count": 1
				},
				{
					"value": "value11",
					"count": 1
				},
				{
					"value": "value13",
					"count": 1
				},
				{
					"value": "value25",
					"count": 1
				},
				{
					"value": "value28",
					"count": 1
				},
				{
					"value": "value29",
					"count": 1
				},
				{
					"value": "value36",
					"count": 1
				},
				{
					"value": "value39",
					"count": 1
				},
				{
					"value": "value66",
					"count": 1
				},
				{
					"value": "value75",
					"count": 1
				},
				{
					"value": "value79",
					"count": 1
				},
				{
					"value": "value89",
					"count": 1
				},
				{
					"value": "value92",
					"count": 1
				}
			]
		},
		"label2": {
			"values": [
				{
					"value": "value24",
					"count": 1
				},
				{
					"value": "value27",
					"count": 1
				},
				{
					"value": "value3",
					"count": 1
				},
				{
					"value": "value33",
					"count": 1
				},
				{
					"value": "value36",
					"count": 1
				},
				{
					"value": "value41",
					"count": 1
				},
				{
					"value": "value49",
					"count": 1
				},
				{
					"value": "value58",
					"count": 1
				},
				{
					"value": "value59",
					"count": 1
				},
				{
					"value": "value62",
					"count": 1
				},
				{
					"value": "value63",
					"count": 1
				},
				{
					"value": "value74",
					"count": 1
				}
			]
		},
		"label3": {
			"values": [
				{
					"value": "value12",
					"count": 1
				},
				{
					"value": "value37",
					"count": 1
				},
				{
					"value": "value38",
					"count": 1
				},
				{
					"value": "value45",
					"count": 1
				},
				{
					"value": "value5",
					"count": 1
				},
				{
					"value": "value50",
					"count": 1
				},
				{
					"value": "value52",
					"count": 1
				},
				{
					"value": "value58",
					"count": 1
				},
				{
					"value": "value61",
					"count": 1
				},
				{
					"value": "value70",
					"count": 1
				},
				{
					"value": "value73",
					"count": 1
				},
				{
					"value": "value75",
					"count": 1
				},
				{
					"value": "value86",
					"count": 1
				},
				{
					"value": "value93",
					"count": 1
				},
				{
					"value": "value97",
					"count": 1
				}
			]
		},
		"label4": {
			"values": [
				{
					"value": "value13",
					"count": 1
				},
				{
					"value": "value41",
					"count": 1
				},
				{
					"value": "value42",
					"count": 2
				},
				{
					"value": "value64",
					"count": 1
				},
				{
					"value": "value75",
					"count": 1
				},
				{
					"value": "value8",
					"count": 2
				},
				{
					"value": "value80",
					"count": 1
				},
				{
					"value": "value83",
					"count": 1
				},
				{
					"value": "value85",
					"count": 1
				},
				{
					"value": "value93",
					"count": 1
				}
			]
		},
		"label5": {
			"values": [
				{
					"value": "value21",
					"count": 1
				},
				{
					"value": "value23",
					"count": 1
				},
				{
					"value": "value28",
					"count": 2
				},
				{
					"value": "value3",
					"count": 1
				},
				{
					"value": "value37",
					"count": 1
				},
				{
					"value": "value39",
					"count": 1
				},
				{
					"value": "value48",
					"count": 1
				},
				{
					"value": "value52",
					"count": 1
				},
				{
					"value": "value56",
					"count": 1
				},
				{
					"value": "value61",
					"count": 1
				},
				{
					"value": "value83",
					"count": 1
				},
				{
					"value": "value9",
					"count": 1
				},
				{
					"value": "value90",
					"count": 1
				},
				{
					"value": "value93",
					"count": 1
				}
			]
		},
		"label6": {
			"values": [
				{
					"value": "value13",
					"count": 1
				},
				{
					"value": "value22",
					"count": 1
				},
				{
					"value": "value23",
					"count": 1
				},
				{
					"value": "value26",
					"count": 1
				},
				{
					"value": "value49",
					"count": 1
				},
				{
					"value": "value5",
					"count": 1
				},
				{
					"value": "value51",
					"count": 1
				},
				{
					"value": "value58",
					"count": 1
				},
				{
					"value": "value61",
					"count": 1
				},
				{
					"value": "value70",
					"count": 1
				},
				{
					"value": "value9",
					"count": 1
				},
				{
					"value": "value90",
					"count": 1
				},
				{
					"value": "value98",
					"count": 1
				}
			]
		},
		"label7": {
			"values": [
				{
					"value": "value0",
					"count": 1
				},
				{
					"value": "value10",
					"count": 1
				},
				{
					"value": "value11",
					"count": 1
				},
				{
					"value": "value14",
					"count": 1
				},
				{
					"value": "value16",
					"count": 1
				},
				{
					"value": "value20",
					"count": 1
				},
				{
					"value": "value23",
					"count": 2
				},
				{
					"value": "value26",
					"count": 1
				},
				{
					"value": "value27",
					"count": 1
				},
				{
					"value": "value32",
					"count": 1
				},
				{
					"value": "value35",
					"count": 1
				},
				{
					"value": "value40",
					"count": 1
				},
				{
					"value": "value55",
					"count": 1
				},
				{
					"value": "value60",
					"count": 1
				},
				{
					"value": "value61",
					"count": 1
				},
				{
					"value": "value64",
					"count": 1
				},
				{
					"value": "value72",
					"count": 1
				},
				{
					"value": "value74",
					"count": 1
				},
				{
					"value": "value79",
					"count": 1
				},
				{
					"value": "value86",
					"count": 2
				},
				{
					"value": "value9",
					"count": 1
				}
			]
		},
		"label8": {
			"values": [
				{
					"value": "value25",
					"count": 1
				},
				{
					"value": "value28",
					"count": 1
				},
				{
					"value": "value32",
					"count": 1
				},
				{
					"value": "value37",
					"count": 1
				},
				{
					"value": "value38",
					"count": 1
				},
				{
					"value": "value42",
					"count": 2
				},
				{
					"value": "value45",
					"count": 1
				},
				{
					"value": "value62",
					"count": 1
				},
				{
					"value": "value7",
					"count": 1
				},
				{
					"value": "value70",
					"count": 1
				},
				{
					"value": "value71",
					"count": 1
				},
				{
					"value": "value81",
					"count": 1
				},
				{
					"value": "value99",
					"count": 1
				}
			]
		},
		"label9": {
			"values": [
				{
					"value": "value11",
					"count": 1
				},
				{
					"value": "value16",
					"count": 1
				},
				{
					"value": "value26",
					"count": 1
				},
				{
					"value": "value29",
					"count": 1
				},
				{
					"value": "value35",
					"count": 1
				},
				{
					"value": "value40",
					"count": 1
				},
				{
					"value": "value41",
					"count": 1
				},
				{
					"value": "value43",
					"count": 1
				},
				{
					"value": "value46",
					"count": 2
				},
				{
					"value": "value55",
					"count": 1
				},
				{
					"value": "value58",
					"count": 1
				},
				{
					"value": "value65",
					"count": 1
				},
				{
					"value": "value71",
					"count": 1
				},
				{
					"value": "value76",
					"count": 1
				},
				{
					"value": "value83",
					"count": 1
				},
				{
					"value": "value95",
					"count": 1
				}
			]
		}
	}
}
`

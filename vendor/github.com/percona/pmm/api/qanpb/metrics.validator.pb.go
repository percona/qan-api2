// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: qanpb/metrics.proto

package qanpb

import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/golang/protobuf/ptypes/timestamp"
import _ "google.golang.org/genproto/googleapis/api/annotations"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *MetricsRequest) Validate() error {
	if this.PeriodStartFrom != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PeriodStartFrom); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PeriodStartFrom", err)
		}
	}
	if this.PeriodStartTo != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PeriodStartTo); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PeriodStartTo", err)
		}
	}
	for _, item := range this.Labels {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Labels", err)
			}
		}
	}
	return nil
}
func (this *MapFieldEntry) Validate() error {
	return nil
}
func (this *MetricsReply) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *MetricValues) Validate() error {
	return nil
}

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: qanpb/profile.proto

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

func (this *ReportRequest) Validate() error {
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
func (this *ReportMapFieldEntry) Validate() error {
	return nil
}
func (this *ReportReply) Validate() error {
	for _, item := range this.Rows {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Rows", err)
			}
		}
	}
	return nil
}
func (this *Row) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	for _, item := range this.Sparkline {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Sparkline", err)
			}
		}
	}
	return nil
}
func (this *Metric) Validate() error {
	if this.Stats != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Stats); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Stats", err)
		}
	}
	return nil
}
func (this *Stat) Validate() error {
	return nil
}
func (this *Point) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}

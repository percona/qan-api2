// Code generated by protoc-gen-go. DO NOT EDIT.
// source: qanpb/profile.proto

package qanpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// ReportRequest defines filtering of metrics report for db server or other dimentions.
type ReportRequest struct {
	PeriodStartFrom      *timestamp.Timestamp   `protobuf:"bytes,1,opt,name=period_start_from,json=periodStartFrom,proto3" json:"period_start_from,omitempty"`
	PeriodStartTo        *timestamp.Timestamp   `protobuf:"bytes,2,opt,name=period_start_to,json=periodStartTo,proto3" json:"period_start_to,omitempty"`
	GroupBy              string                 `protobuf:"bytes,3,opt,name=group_by,json=groupBy,proto3" json:"group_by,omitempty"`
	Labels               []*ReportMapFieldEntry `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty"`
	Columns              []string               `protobuf:"bytes,5,rep,name=columns,proto3" json:"columns,omitempty"`
	OrderBy              string                 `protobuf:"bytes,6,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Offset               uint32                 `protobuf:"varint,7,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                uint32                 `protobuf:"varint,8,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ReportRequest) Reset()         { *m = ReportRequest{} }
func (m *ReportRequest) String() string { return proto.CompactTextString(m) }
func (*ReportRequest) ProtoMessage()    {}
func (*ReportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_81344f06bfb65433, []int{0}
}
func (m *ReportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportRequest.Unmarshal(m, b)
}
func (m *ReportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportRequest.Marshal(b, m, deterministic)
}
func (dst *ReportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportRequest.Merge(dst, src)
}
func (m *ReportRequest) XXX_Size() int {
	return xxx_messageInfo_ReportRequest.Size(m)
}
func (m *ReportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReportRequest proto.InternalMessageInfo

func (m *ReportRequest) GetPeriodStartFrom() *timestamp.Timestamp {
	if m != nil {
		return m.PeriodStartFrom
	}
	return nil
}

func (m *ReportRequest) GetPeriodStartTo() *timestamp.Timestamp {
	if m != nil {
		return m.PeriodStartTo
	}
	return nil
}

func (m *ReportRequest) GetGroupBy() string {
	if m != nil {
		return m.GroupBy
	}
	return ""
}

func (m *ReportRequest) GetLabels() []*ReportMapFieldEntry {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *ReportRequest) GetColumns() []string {
	if m != nil {
		return m.Columns
	}
	return nil
}

func (m *ReportRequest) GetOrderBy() string {
	if m != nil {
		return m.OrderBy
	}
	return ""
}

func (m *ReportRequest) GetOffset() uint32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *ReportRequest) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

// ReportMapFieldEntry allows to pass labels/dimentions in form like {"d_server": ["db1", "db2"...]}.
type ReportMapFieldEntry struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []string `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportMapFieldEntry) Reset()         { *m = ReportMapFieldEntry{} }
func (m *ReportMapFieldEntry) String() string { return proto.CompactTextString(m) }
func (*ReportMapFieldEntry) ProtoMessage()    {}
func (*ReportMapFieldEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_81344f06bfb65433, []int{1}
}
func (m *ReportMapFieldEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportMapFieldEntry.Unmarshal(m, b)
}
func (m *ReportMapFieldEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportMapFieldEntry.Marshal(b, m, deterministic)
}
func (dst *ReportMapFieldEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportMapFieldEntry.Merge(dst, src)
}
func (m *ReportMapFieldEntry) XXX_Size() int {
	return xxx_messageInfo_ReportMapFieldEntry.Size(m)
}
func (m *ReportMapFieldEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportMapFieldEntry.DiscardUnknown(m)
}

var xxx_messageInfo_ReportMapFieldEntry proto.InternalMessageInfo

func (m *ReportMapFieldEntry) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ReportMapFieldEntry) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

// ReportReply is list of reports per quieryids, hosts etc.
type ReportReply struct {
	TotalRows            uint32   `protobuf:"varint,1,opt,name=total_rows,json=totalRows,proto3" json:"total_rows,omitempty"`
	Offset               uint32   `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                uint32   `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Rows                 []*Row   `protobuf:"bytes,4,rep,name=rows,proto3" json:"rows,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportReply) Reset()         { *m = ReportReply{} }
func (m *ReportReply) String() string { return proto.CompactTextString(m) }
func (*ReportReply) ProtoMessage()    {}
func (*ReportReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_81344f06bfb65433, []int{2}
}
func (m *ReportReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportReply.Unmarshal(m, b)
}
func (m *ReportReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportReply.Marshal(b, m, deterministic)
}
func (dst *ReportReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportReply.Merge(dst, src)
}
func (m *ReportReply) XXX_Size() int {
	return xxx_messageInfo_ReportReply.Size(m)
}
func (m *ReportReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportReply.DiscardUnknown(m)
}

var xxx_messageInfo_ReportReply proto.InternalMessageInfo

func (m *ReportReply) GetTotalRows() uint32 {
	if m != nil {
		return m.TotalRows
	}
	return 0
}

func (m *ReportReply) GetOffset() uint32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *ReportReply) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ReportReply) GetRows() []*Row {
	if m != nil {
		return m.Rows
	}
	return nil
}

// Row define metrics for selected dimention.
type Row struct {
	Rank                 uint32             `protobuf:"varint,1,opt,name=rank,proto3" json:"rank,omitempty"`
	Dimension            string             `protobuf:"bytes,2,opt,name=dimension,proto3" json:"dimension,omitempty"`
	Metrics              map[string]*Metric `protobuf:"bytes,3,rep,name=metrics,proto3" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Sparkline            []*Point           `protobuf:"bytes,4,rep,name=sparkline,proto3" json:"sparkline,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Row) Reset()         { *m = Row{} }
func (m *Row) String() string { return proto.CompactTextString(m) }
func (*Row) ProtoMessage()    {}
func (*Row) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_81344f06bfb65433, []int{3}
}
func (m *Row) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Row.Unmarshal(m, b)
}
func (m *Row) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Row.Marshal(b, m, deterministic)
}
func (dst *Row) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Row.Merge(dst, src)
}
func (m *Row) XXX_Size() int {
	return xxx_messageInfo_Row.Size(m)
}
func (m *Row) XXX_DiscardUnknown() {
	xxx_messageInfo_Row.DiscardUnknown(m)
}

var xxx_messageInfo_Row proto.InternalMessageInfo

func (m *Row) GetRank() uint32 {
	if m != nil {
		return m.Rank
	}
	return 0
}

func (m *Row) GetDimension() string {
	if m != nil {
		return m.Dimension
	}
	return ""
}

func (m *Row) GetMetrics() map[string]*Metric {
	if m != nil {
		return m.Metrics
	}
	return nil
}

func (m *Row) GetSparkline() []*Point {
	if m != nil {
		return m.Sparkline
	}
	return nil
}

// Metric cell.
type Metric struct {
	Stats                *Stat    `protobuf:"bytes,1,opt,name=stats,proto3" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_81344f06bfb65433, []int{4}
}
func (m *Metric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metric.Unmarshal(m, b)
}
func (m *Metric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metric.Marshal(b, m, deterministic)
}
func (dst *Metric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metric.Merge(dst, src)
}
func (m *Metric) XXX_Size() int {
	return xxx_messageInfo_Metric.Size(m)
}
func (m *Metric) XXX_DiscardUnknown() {
	xxx_messageInfo_Metric.DiscardUnknown(m)
}

var xxx_messageInfo_Metric proto.InternalMessageInfo

func (m *Metric) GetStats() *Stat {
	if m != nil {
		return m.Stats
	}
	return nil
}

// Stat is statistics of specific metric.
type Stat struct {
	Rate                 float32  `protobuf:"fixed32,1,opt,name=rate,proto3" json:"rate,omitempty"`
	Cnt                  float32  `protobuf:"fixed32,2,opt,name=cnt,proto3" json:"cnt,omitempty"`
	Sum                  float32  `protobuf:"fixed32,3,opt,name=sum,proto3" json:"sum,omitempty"`
	Min                  float32  `protobuf:"fixed32,4,opt,name=min,proto3" json:"min,omitempty"`
	Max                  float32  `protobuf:"fixed32,5,opt,name=max,proto3" json:"max,omitempty"`
	P99                  float32  `protobuf:"fixed32,6,opt,name=p99,proto3" json:"p99,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stat) Reset()         { *m = Stat{} }
func (m *Stat) String() string { return proto.CompactTextString(m) }
func (*Stat) ProtoMessage()    {}
func (*Stat) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_81344f06bfb65433, []int{5}
}
func (m *Stat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stat.Unmarshal(m, b)
}
func (m *Stat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stat.Marshal(b, m, deterministic)
}
func (dst *Stat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stat.Merge(dst, src)
}
func (m *Stat) XXX_Size() int {
	return xxx_messageInfo_Stat.Size(m)
}
func (m *Stat) XXX_DiscardUnknown() {
	xxx_messageInfo_Stat.DiscardUnknown(m)
}

var xxx_messageInfo_Stat proto.InternalMessageInfo

func (m *Stat) GetRate() float32 {
	if m != nil {
		return m.Rate
	}
	return 0
}

func (m *Stat) GetCnt() float32 {
	if m != nil {
		return m.Cnt
	}
	return 0
}

func (m *Stat) GetSum() float32 {
	if m != nil {
		return m.Sum
	}
	return 0
}

func (m *Stat) GetMin() float32 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *Stat) GetMax() float32 {
	if m != nil {
		return m.Max
	}
	return 0
}

func (m *Stat) GetP99() float32 {
	if m != nil {
		return m.P99
	}
	return 0
}

// Point contains values that represents abscissa (time) and ordinate (volume etc.)
// of every point in a coordinate system of Sparklines.
type Point struct {
	Values               map[string]float32 `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Point) Reset()         { *m = Point{} }
func (m *Point) String() string { return proto.CompactTextString(m) }
func (*Point) ProtoMessage()    {}
func (*Point) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_81344f06bfb65433, []int{6}
}
func (m *Point) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Point.Unmarshal(m, b)
}
func (m *Point) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Point.Marshal(b, m, deterministic)
}
func (dst *Point) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Point.Merge(dst, src)
}
func (m *Point) XXX_Size() int {
	return xxx_messageInfo_Point.Size(m)
}
func (m *Point) XXX_DiscardUnknown() {
	xxx_messageInfo_Point.DiscardUnknown(m)
}

var xxx_messageInfo_Point proto.InternalMessageInfo

func (m *Point) GetValues() map[string]float32 {
	if m != nil {
		return m.Values
	}
	return nil
}

func init() {
	proto.RegisterType((*ReportRequest)(nil), "qan.ReportRequest")
	proto.RegisterType((*ReportMapFieldEntry)(nil), "qan.ReportMapFieldEntry")
	proto.RegisterType((*ReportReply)(nil), "qan.ReportReply")
	proto.RegisterType((*Row)(nil), "qan.Row")
	proto.RegisterMapType((map[string]*Metric)(nil), "qan.Row.MetricsEntry")
	proto.RegisterType((*Metric)(nil), "qan.Metric")
	proto.RegisterType((*Stat)(nil), "qan.Stat")
	proto.RegisterType((*Point)(nil), "qan.Point")
	proto.RegisterMapType((map[string]float32)(nil), "qan.Point.ValuesEntry")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProfileClient is the client API for Profile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProfileClient interface {
	// GetReport returns list of metrics group by queryid or other dimentions.
	GetReport(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportReply, error)
}

type profileClient struct {
	cc *grpc.ClientConn
}

func NewProfileClient(cc *grpc.ClientConn) ProfileClient {
	return &profileClient{cc}
}

func (c *profileClient) GetReport(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportReply, error) {
	out := new(ReportReply)
	err := c.cc.Invoke(ctx, "/qan.Profile/GetReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServer is the server API for Profile service.
type ProfileServer interface {
	// GetReport returns list of metrics group by queryid or other dimentions.
	GetReport(context.Context, *ReportRequest) (*ReportReply, error)
}

func RegisterProfileServer(s *grpc.Server, srv ProfileServer) {
	s.RegisterService(&_Profile_serviceDesc, srv)
}

func _Profile_GetReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/qan.Profile/GetReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetReport(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Profile_serviceDesc = grpc.ServiceDesc{
	ServiceName: "qan.Profile",
	HandlerType: (*ProfileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReport",
			Handler:    _Profile_GetReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "qanpb/profile.proto",
}

func init() { proto.RegisterFile("qanpb/profile.proto", fileDescriptor_profile_81344f06bfb65433) }

var fileDescriptor_profile_81344f06bfb65433 = []byte{
	// 719 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0x4f, 0x6f, 0xd3, 0x4e,
	0x10, 0x55, 0xec, 0xfc, 0xa9, 0x27, 0xbf, 0xa8, 0xed, 0xf6, 0x47, 0x65, 0xa2, 0xa0, 0x1a, 0xc3,
	0x21, 0x20, 0x6a, 0x43, 0xb9, 0xd0, 0x4a, 0x1c, 0x1a, 0xa9, 0xed, 0xa9, 0xa2, 0x6c, 0x2b, 0x0e,
	0xbd, 0x44, 0x9b, 0x64, 0x13, 0xad, 0x6a, 0xef, 0x3a, 0xbb, 0x9b, 0xa6, 0xb9, 0xf2, 0x11, 0xe0,
	0x63, 0x71, 0xe4, 0xc8, 0x95, 0x0f, 0x82, 0x76, 0xd7, 0x21, 0xa9, 0xa8, 0xe0, 0x94, 0x99, 0xe7,
	0x99, 0xf7, 0xc6, 0xcf, 0x93, 0x81, 0x9d, 0x29, 0xe1, 0xc5, 0x20, 0x2d, 0xa4, 0x18, 0xb3, 0x8c,
	0x26, 0x85, 0x14, 0x5a, 0x20, 0x7f, 0x4a, 0x78, 0xbb, 0x33, 0x11, 0x62, 0x92, 0xd1, 0x94, 0x14,
	0x2c, 0x25, 0x9c, 0x0b, 0x4d, 0x34, 0x13, 0x5c, 0xb9, 0x92, 0xf6, 0x5e, 0xf9, 0xd4, 0x66, 0x83,
	0xd9, 0x38, 0xd5, 0x2c, 0xa7, 0x4a, 0x93, 0xbc, 0x28, 0x0b, 0x5e, 0xd9, 0x9f, 0xe1, 0xfe, 0x84,
	0xf2, 0x7d, 0x35, 0x27, 0x93, 0x09, 0x95, 0xa9, 0x28, 0x2c, 0xc5, 0x9f, 0x74, 0xf1, 0x37, 0x0f,
	0x5a, 0x98, 0x16, 0x42, 0x6a, 0x4c, 0xa7, 0x33, 0xaa, 0x34, 0x3a, 0x85, 0xed, 0x82, 0x4a, 0x26,
	0x46, 0x7d, 0xa5, 0x89, 0xd4, 0xfd, 0xb1, 0x14, 0x79, 0x58, 0x89, 0x2a, 0xdd, 0xe6, 0x41, 0x3b,
	0x71, 0xe2, 0xc9, 0x52, 0x3c, 0xb9, 0x5a, 0x8a, 0xe3, 0x4d, 0xd7, 0x74, 0x69, 0x7a, 0x4e, 0xa5,
	0xc8, 0x51, 0x0f, 0x36, 0xef, 0xf1, 0x68, 0x11, 0x7a, 0xff, 0x64, 0x69, 0xad, 0xb1, 0x5c, 0x09,
	0xf4, 0x18, 0x36, 0x26, 0x52, 0xcc, 0x8a, 0xfe, 0x60, 0x11, 0xfa, 0x51, 0xa5, 0x1b, 0xe0, 0x86,
	0xcd, 0x7b, 0x0b, 0xf4, 0x1a, 0xea, 0x19, 0x19, 0xd0, 0x4c, 0x85, 0xd5, 0xc8, 0xef, 0x36, 0x0f,
	0xc2, 0x64, 0x4a, 0x78, 0xe2, 0x5e, 0xe5, 0x9c, 0x14, 0xa7, 0x8c, 0x66, 0xa3, 0x13, 0xae, 0xe5,
	0x02, 0x97, 0x75, 0x28, 0x84, 0xc6, 0x50, 0x64, 0xb3, 0x9c, 0xab, 0xb0, 0x16, 0xf9, 0x86, 0xab,
	0x4c, 0x8d, 0x8c, 0x90, 0x23, 0x2a, 0x8d, 0x4c, 0xdd, 0xc9, 0xd8, 0xbc, 0xb7, 0x40, 0xbb, 0x50,
	0x17, 0xe3, 0xb1, 0xa2, 0x3a, 0x6c, 0x44, 0x95, 0x6e, 0x0b, 0x97, 0x19, 0xfa, 0x1f, 0x6a, 0x19,
	0xcb, 0x99, 0x0e, 0x37, 0x2c, 0xec, 0x92, 0xf8, 0x3d, 0xec, 0x3c, 0x30, 0x01, 0xda, 0x02, 0xff,
	0x86, 0x2e, 0xac, 0x89, 0x01, 0x36, 0xa1, 0x69, 0xbf, 0x25, 0xd9, 0x8c, 0x86, 0x9e, 0x9d, 0xc4,
	0x25, 0xf1, 0x1d, 0x34, 0x97, 0xdf, 0xa2, 0xc8, 0x16, 0xe8, 0x09, 0x80, 0x16, 0x9a, 0x64, 0x7d,
	0x29, 0xe6, 0xca, 0x76, 0xb7, 0x70, 0x60, 0x11, 0x2c, 0xe6, 0x6a, 0x6d, 0x34, 0xef, 0xe1, 0xd1,
	0xfc, 0xb5, 0xd1, 0x50, 0x07, 0xaa, 0x96, 0xc6, 0xb9, 0xb5, 0xe1, 0xdc, 0x12, 0x73, 0x6c, 0xd1,
	0xf8, 0x47, 0x05, 0x7c, 0x2c, 0xe6, 0x08, 0x41, 0x55, 0x12, 0x7e, 0x53, 0x8a, 0xd9, 0x18, 0x75,
	0x20, 0x18, 0xb1, 0x9c, 0x72, 0xc5, 0x04, 0xb7, 0x52, 0x01, 0x5e, 0x01, 0x28, 0x85, 0x46, 0x4e,
	0xb5, 0x64, 0x43, 0x15, 0xfa, 0x96, 0xfa, 0xd1, 0x92, 0x3a, 0x39, 0x77, 0xb8, 0xfb, 0x0a, 0xcb,
	0x2a, 0xd4, 0x85, 0x40, 0x15, 0x44, 0xde, 0x64, 0x8c, 0xd3, 0x72, 0x1a, 0xb0, 0x2d, 0x17, 0x82,
	0x71, 0x8d, 0x57, 0x0f, 0xdb, 0x67, 0xf0, 0xdf, 0x3a, 0xc5, 0x03, 0x36, 0x3e, 0x5d, 0xd9, 0x68,
	0x36, 0xab, 0x69, 0x79, 0x5c, 0x4f, 0xe9, 0xe9, 0x91, 0xf7, 0xae, 0x12, 0xbf, 0x80, 0xba, 0x03,
	0xd1, 0x1e, 0xd4, 0x94, 0x26, 0x5a, 0x95, 0x0b, 0x1d, 0xd8, 0x86, 0x4b, 0x4d, 0x34, 0x76, 0x78,
	0xac, 0xa1, 0x6a, 0x52, 0x67, 0x84, 0xa6, 0xb6, 0xce, 0xc3, 0x36, 0x36, 0xfa, 0x43, 0xee, 0xdc,
	0xf6, 0xb0, 0x09, 0x0d, 0xa2, 0x66, 0xb9, 0x35, 0xda, 0xc3, 0x26, 0x34, 0x48, 0xce, 0x78, 0x58,
	0x75, 0x48, 0xce, 0xb8, 0x45, 0xc8, 0x5d, 0x58, 0x2b, 0x11, 0x72, 0x67, 0x90, 0xe2, 0xf0, 0xd0,
	0x6e, 0x9a, 0x87, 0x4d, 0x18, 0x4b, 0xa8, 0xd9, 0xb7, 0x47, 0x09, 0xd4, 0xed, 0xd8, 0x66, 0x40,
	0xe3, 0xcc, 0xee, 0xca, 0x99, 0xe4, 0x93, 0x7d, 0x50, 0xee, 0xb4, 0xab, 0x6a, 0x1f, 0x42, 0x73,
	0x0d, 0xfe, 0xfb, 0xa2, 0x19, 0xb5, 0x95, 0x29, 0x07, 0xd7, 0xd0, 0xb8, 0x70, 0xc7, 0x07, 0x7d,
	0x80, 0xe0, 0x8c, 0x6a, 0xb7, 0x7a, 0x08, 0xad, 0xfd, 0x91, 0xca, 0x9b, 0xd0, 0xde, 0xba, 0x87,
	0x15, 0xd9, 0x22, 0xee, 0x7c, 0xfe, 0xfe, 0xf3, 0xab, 0xb7, 0x1b, 0x6f, 0xa7, 0xb7, 0x6f, 0xd2,
	0x29, 0xe1, 0xe9, 0x6f, 0x82, 0xa3, 0xca, 0xcb, 0xde, 0xc7, 0x2f, 0xc7, 0x67, 0xf8, 0x04, 0x1a,
	0x23, 0x3a, 0x26, 0xb3, 0x4c, 0xa3, 0x23, 0x40, 0xc7, 0x3c, 0xa2, 0x52, 0x0a, 0x19, 0x49, 0xaa,
	0x0a, 0xc1, 0x15, 0x4d, 0xd0, 0x73, 0x88, 0xdb, 0xd1, 0xb3, 0x74, 0x44, 0xc7, 0x8c, 0x33, 0x77,
	0xa0, 0xec, 0x51, 0x3c, 0x31, 0x75, 0xb8, 0x2c, 0xbb, 0xae, 0x59, 0x6c, 0x50, 0xb7, 0xd7, 0xe2,
	0xed, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xda, 0x3a, 0x01, 0x70, 0x38, 0x05, 0x00, 0x00,
}

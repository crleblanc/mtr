// Code generated by protoc-gen-go.
// source: data.proto
// DO NOT EDIT!

/*
Package mtrpb is a generated protocol buffer package.

It is generated from these files:
	data.proto
	field.proto
	tag.proto

It has these top-level messages:
	DataLatencySummary
	DataLatencySummaryResult
	DataSite
	DataSiteResult
	DataLatencyTag
	DataLatencyTagResult
	FieldMetricSummary
	FieldMetricSummaryResult
	FieldMetricTag
	FieldMetricTagResult
	FieldMetricThreshold
	FieldMetricThresholdResult
	Tag
	TagResult
	TagSearchResult
*/
package mtrpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

// DataLatencySummary is a summary of data latency metrics for each site.
// mean should not be 0.  fifty and ninety may be unknown (0).
// If upper == lower == 0 then no threshold has been set on the metric.
type DataLatencySummary struct {
	// The siteID for the metric e.g., TAUP
	SiteID string `protobuf:"bytes,1,opt,name=site_iD,json=siteID" json:"site_iD,omitempty"`
	// The typeID for the metric e.g., latency.strong
	TypeID string `protobuf:"bytes,2,opt,name=type_iD,json=typeID" json:"type_iD,omitempty"`
	// Unix time in seconds for the metric value (don't need nanos).
	Seconds int64 `protobuf:"varint,3,opt,name=seconds" json:"seconds,omitempty"`
	// The mean latency
	Mean int32 `protobuf:"varint,4,opt,name=mean" json:"mean,omitempty"`
	// The fiftieth percentile value.  Might be unknown (0)
	Fifty int32 `protobuf:"varint,5,opt,name=fifty" json:"fifty,omitempty"`
	// The ninetieth percentile value.  Might be unknown (0)
	Ninety int32 `protobuf:"varint,6,opt,name=ninety" json:"ninety,omitempty"`
	// The upper threshold for the metric to be good.
	Upper int32 `protobuf:"varint,7,opt,name=upper" json:"upper,omitempty"`
	// The lower threshold for the metric to be good.
	Lower int32 `protobuf:"varint,8,opt,name=lower" json:"lower,omitempty"`
}

func (m *DataLatencySummary) Reset()                    { *m = DataLatencySummary{} }
func (m *DataLatencySummary) String() string            { return proto.CompactTextString(m) }
func (*DataLatencySummary) ProtoMessage()               {}
func (*DataLatencySummary) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type DataLatencySummaryResult struct {
	Result []*DataLatencySummary `protobuf:"bytes,1,rep,name=result" json:"result,omitempty"`
}

func (m *DataLatencySummaryResult) Reset()                    { *m = DataLatencySummaryResult{} }
func (m *DataLatencySummaryResult) String() string            { return proto.CompactTextString(m) }
func (*DataLatencySummaryResult) ProtoMessage()               {}
func (*DataLatencySummaryResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DataLatencySummaryResult) GetResult() []*DataLatencySummary {
	if m != nil {
		return m.Result
	}
	return nil
}

type DataSite struct {
	// The siteID for the metric e.g., TAUP
	SiteID string `protobuf:"bytes,1,opt,name=site_iD,json=siteID" json:"site_iD,omitempty"`
	// The site latitude - not usually accurate enough for meta data
	Latitude float64 `protobuf:"fixed64,2,opt,name=latitude" json:"latitude,omitempty"`
	// The site longitude - not usually accurate enough for meta data
	Longitude float64 `protobuf:"fixed64,3,opt,name=longitude" json:"longitude,omitempty"`
}

func (m *DataSite) Reset()                    { *m = DataSite{} }
func (m *DataSite) String() string            { return proto.CompactTextString(m) }
func (*DataSite) ProtoMessage()               {}
func (*DataSite) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type DataSiteResult struct {
	Result []*DataSite `protobuf:"bytes,1,rep,name=result" json:"result,omitempty"`
}

func (m *DataSiteResult) Reset()                    { *m = DataSiteResult{} }
func (m *DataSiteResult) String() string            { return proto.CompactTextString(m) }
func (*DataSiteResult) ProtoMessage()               {}
func (*DataSiteResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DataSiteResult) GetResult() []*DataSite {
	if m != nil {
		return m.Result
	}
	return nil
}

type DataLatencyTag struct {
	// The siteID for the latency e.g., TAUP
	SiteID string `protobuf:"bytes,1,opt,name=site_iD,json=siteID" json:"site_iD,omitempty"`
	// The typeID for the latency e.g., latency.gnss.1hz
	TypeID string `protobuf:"bytes,2,opt,name=type_iD,json=typeID" json:"type_iD,omitempty"`
	// The tag for the latency e.g., TAUP
	Tag string `protobuf:"bytes,3,opt,name=tag" json:"tag,omitempty"`
}

func (m *DataLatencyTag) Reset()                    { *m = DataLatencyTag{} }
func (m *DataLatencyTag) String() string            { return proto.CompactTextString(m) }
func (*DataLatencyTag) ProtoMessage()               {}
func (*DataLatencyTag) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type DataLatencyTagResult struct {
	Result []*DataLatencyTag `protobuf:"bytes,1,rep,name=result" json:"result,omitempty"`
}

func (m *DataLatencyTagResult) Reset()                    { *m = DataLatencyTagResult{} }
func (m *DataLatencyTagResult) String() string            { return proto.CompactTextString(m) }
func (*DataLatencyTagResult) ProtoMessage()               {}
func (*DataLatencyTagResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *DataLatencyTagResult) GetResult() []*DataLatencyTag {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*DataLatencySummary)(nil), "mtrpb.DataLatencySummary")
	proto.RegisterType((*DataLatencySummaryResult)(nil), "mtrpb.DataLatencySummaryResult")
	proto.RegisterType((*DataSite)(nil), "mtrpb.DataSite")
	proto.RegisterType((*DataSiteResult)(nil), "mtrpb.DataSiteResult")
	proto.RegisterType((*DataLatencyTag)(nil), "mtrpb.DataLatencyTag")
	proto.RegisterType((*DataLatencyTagResult)(nil), "mtrpb.DataLatencyTagResult")
}

var fileDescriptor0 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x52, 0x4f, 0x4b, 0xc3, 0x30,
	0x14, 0xa7, 0x76, 0xed, 0xb6, 0x27, 0xa8, 0x84, 0xa9, 0x51, 0x3c, 0x48, 0x2f, 0xee, 0xe2, 0x40,
	0x3d, 0x79, 0x95, 0x7a, 0x10, 0xf4, 0x92, 0xed, 0x24, 0x88, 0x64, 0x6b, 0x56, 0x0a, 0x6d, 0x5a,
	0xda, 0x14, 0xe9, 0x97, 0xf4, 0x33, 0x99, 0xbc, 0xa4, 0xfe, 0xd9, 0xf4, 0xe2, 0xed, 0xfd, 0xfe,
	0xbc, 0xf0, 0x7b, 0x2f, 0x0f, 0x20, 0xe1, 0x8a, 0xcf, 0xaa, 0xba, 0x54, 0x25, 0x09, 0x0a, 0x55,
	0x57, 0xcb, 0xe8, 0xdd, 0x03, 0x12, 0x6b, 0xf6, 0x91, 0x2b, 0x21, 0x57, 0xdd, 0xbc, 0x2d, 0x0a,
	0x5e, 0x77, 0xe4, 0x18, 0x86, 0x4d, 0xa6, 0xc4, 0x6b, 0x16, 0x53, 0xef, 0xdc, 0x9b, 0x8e, 0x59,
	0x68, 0xe0, 0x43, 0x6c, 0x04, 0xd5, 0x55, 0x28, 0xec, 0x58, 0xc1, 0x40, 0x2d, 0x50, 0xdd, 0x21,
	0x56, 0xa5, 0x4c, 0x1a, 0xea, 0x6b, 0xc1, 0x67, 0x3d, 0x24, 0x04, 0x06, 0x85, 0xe0, 0x92, 0x0e,
	0x34, 0x1d, 0x30, 0xac, 0xc9, 0x04, 0x82, 0x75, 0xb6, 0x56, 0x1d, 0x0d, 0x90, 0xb4, 0x80, 0x1c,
	0x41, 0x28, 0x33, 0x29, 0x34, 0x1d, 0x22, 0xed, 0x90, 0x71, 0xb7, 0x55, 0x25, 0x6a, 0x3a, 0xb4,
	0x6e, 0x04, 0x86, 0xcd, 0xcb, 0x37, 0xcd, 0x8e, 0x2c, 0x8b, 0x20, 0x7a, 0x02, 0xba, 0x3d, 0x0f,
	0x13, 0x4d, 0x9b, 0x2b, 0x72, 0x05, 0x61, 0x8d, 0x95, 0x1e, 0xca, 0x9f, 0xee, 0x5e, 0x9f, 0xcc,
	0x70, 0x09, 0xb3, 0x5f, 0x1a, 0x9c, 0x31, 0x7a, 0x81, 0x91, 0x51, 0xe7, 0x7a, 0xfa, 0xbf, 0x97,
	0x72, 0x0a, 0xa3, 0x9c, 0xab, 0x4c, 0xb5, 0x89, 0xc0, 0xad, 0x78, 0xec, 0x13, 0x93, 0x33, 0x18,
	0xe7, 0xa5, 0x4c, 0xad, 0xe8, 0xa3, 0xf8, 0x45, 0x44, 0xb7, 0xb0, 0xd7, 0x3f, 0xef, 0x32, 0x5e,
	0x6c, 0x64, 0xdc, 0xff, 0x96, 0x11, 0x6d, 0x7d, 0xb2, 0x85, 0x6d, 0x75, 0xb9, 0x17, 0x3c, 0xfd,
	0xc7, 0xa7, 0x1d, 0x80, 0xaf, 0x78, 0x8a, 0xb1, 0xc6, 0xcc, 0x94, 0xd1, 0x3d, 0x4c, 0x7e, 0xbe,
	0xea, 0x62, 0x5d, 0x6e, 0xc4, 0x3a, 0xdc, 0x5e, 0x9d, 0x31, 0x3b, 0xd3, 0xdd, 0xf0, 0xd9, 0xde,
	0xd7, 0x32, 0xc4, 0x6b, 0xbb, 0xf9, 0x08, 0x00, 0x00, 0xff, 0xff, 0xb6, 0x0b, 0x40, 0x01, 0x7b,
	0x02, 0x00, 0x00,
}

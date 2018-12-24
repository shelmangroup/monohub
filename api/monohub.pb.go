// Code generated by protoc-gen-go. DO NOT EDIT.
// source: monohub.proto

package monohub

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
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

type PingRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_monohub_9e6214d81d0d2c8f, []int{0}
}
func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (dst *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(dst, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

type PingReply struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingReply) Reset()         { *m = PingReply{} }
func (m *PingReply) String() string { return proto.CompactTextString(m) }
func (*PingReply) ProtoMessage()    {}
func (*PingReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_monohub_9e6214d81d0d2c8f, []int{1}
}
func (m *PingReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingReply.Unmarshal(m, b)
}
func (m *PingReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingReply.Marshal(b, m, deterministic)
}
func (dst *PingReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingReply.Merge(dst, src)
}
func (m *PingReply) XXX_Size() int {
	return xxx_messageInfo_PingReply.Size(m)
}
func (m *PingReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PingReply.DiscardUnknown(m)
}

var xxx_messageInfo_PingReply proto.InternalMessageInfo

func (m *PingReply) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*PingRequest)(nil), "monohub.PingRequest")
	proto.RegisterType((*PingReply)(nil), "monohub.PingReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MonohubClient is the client API for Monohub service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MonohubClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error)
}

type monohubClient struct {
	cc *grpc.ClientConn
}

func NewMonohubClient(cc *grpc.ClientConn) MonohubClient {
	return &monohubClient{cc}
}

func (c *monohubClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error) {
	out := new(PingReply)
	err := c.cc.Invoke(ctx, "/monohub.Monohub/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonohubServer is the server API for Monohub service.
type MonohubServer interface {
	Ping(context.Context, *PingRequest) (*PingReply, error)
}

func RegisterMonohubServer(s *grpc.Server, srv MonohubServer) {
	s.RegisterService(&_Monohub_serviceDesc, srv)
}

func _Monohub_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonohubServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/monohub.Monohub/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonohubServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Monohub_serviceDesc = grpc.ServiceDesc{
	ServiceName: "monohub.Monohub",
	HandlerType: (*MonohubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Monohub_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "monohub.proto",
}

func init() { proto.RegisterFile("monohub.proto", fileDescriptor_monohub_9e6214d81d0d2c8f) }

var fileDescriptor_monohub_9e6214d81d0d2c8f = []byte{
	// 163 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0xcd, 0xcf, 0xcb,
	0xcf, 0x28, 0x4d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0xa5, 0x64, 0xd2,
	0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0x13, 0x0b, 0x32, 0xf5, 0x13, 0xf3, 0xf2, 0xf2, 0x4b, 0x12,
	0x4b, 0x32, 0xf3, 0xf3, 0x8a, 0x21, 0xca, 0x94, 0x78, 0xb9, 0xb8, 0x03, 0x32, 0xf3, 0xd2, 0x83,
	0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x94, 0x54, 0xb9, 0x38, 0x21, 0xdc, 0x82, 0x9c, 0x4a, 0x21,
	0x09, 0x2e, 0xf6, 0xb2, 0xd4, 0xa2, 0xe2, 0xcc, 0xfc, 0x3c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce,
	0x20, 0x18, 0xd7, 0xc8, 0x97, 0x8b, 0xdd, 0x17, 0x62, 0xbc, 0x90, 0x13, 0x17, 0x0b, 0x48, 0x87,
	0x90, 0x88, 0x1e, 0xcc, 0x7e, 0x24, 0xf3, 0xa4, 0x84, 0xd0, 0x44, 0x0b, 0x72, 0x2a, 0x95, 0x04,
	0x9a, 0x2e, 0x3f, 0x99, 0xcc, 0xc4, 0x25, 0xc4, 0xa1, 0x5f, 0x66, 0xa8, 0x5f, 0x90, 0x99, 0x97,
	0x9e, 0xc4, 0x06, 0x76, 0x8b, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xe2, 0xff, 0xf3, 0x2f, 0xc3,
	0x00, 0x00, 0x00,
}

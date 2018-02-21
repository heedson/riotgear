// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	api.proto

It has these top-level messages:
	EchoMsg
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type EchoMsg struct {
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *EchoMsg) Reset()                    { *m = EchoMsg{} }
func (m *EchoMsg) String() string            { return proto1.CompactTextString(m) }
func (*EchoMsg) ProtoMessage()               {}
func (*EchoMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EchoMsg) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto1.RegisterType((*EchoMsg)(nil), "EchoMsg")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EchoTest service

type EchoTestClient interface {
	Echo(ctx context.Context, in *EchoMsg, opts ...grpc.CallOption) (*EchoMsg, error)
}

type echoTestClient struct {
	cc *grpc.ClientConn
}

func NewEchoTestClient(cc *grpc.ClientConn) EchoTestClient {
	return &echoTestClient{cc}
}

func (c *echoTestClient) Echo(ctx context.Context, in *EchoMsg, opts ...grpc.CallOption) (*EchoMsg, error) {
	out := new(EchoMsg)
	err := grpc.Invoke(ctx, "/EchoTest/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EchoTest service

type EchoTestServer interface {
	Echo(context.Context, *EchoMsg) (*EchoMsg, error)
}

func RegisterEchoTestServer(s *grpc.Server, srv EchoTestServer) {
	s.RegisterService(&_EchoTest_serviceDesc, srv)
}

func _EchoTest_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoTestServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EchoTest/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoTestServer).Echo(ctx, req.(*EchoMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _EchoTest_serviceDesc = grpc.ServiceDesc{
	ServiceName: "EchoTest",
	HandlerType: (*EchoTestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoTest_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto1.RegisterFile("api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c,
	0xc8, 0xd4, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0x86, 0xc8, 0x2a,
	0xc9, 0x73, 0xb1, 0xbb, 0x26, 0x67, 0xe4, 0xfb, 0x16, 0xa7, 0x0b, 0x89, 0x70, 0xb1, 0x96, 0x25,
	0xe6, 0x94, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x38, 0x46, 0xf6, 0x5c, 0x1c,
	0x20, 0x05, 0x21, 0xa9, 0xc5, 0x25, 0x42, 0xc6, 0x5c, 0x2c, 0x20, 0xb6, 0x10, 0x87, 0x1e, 0x54,
	0x8f, 0x14, 0x9c, 0xa5, 0x24, 0xde, 0x74, 0xf9, 0xc9, 0x64, 0x26, 0x41, 0x25, 0x1e, 0xb0, 0x45,
	0x65, 0x86, 0xfa, 0xa9, 0xc9, 0x19, 0xf9, 0x56, 0x8c, 0x5a, 0x4e, 0xca, 0x51, 0x8a, 0xe9, 0x99,
	0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x19, 0xa9, 0xa9, 0x29, 0xc5, 0xf9, 0x79,
	0xfa, 0x45, 0x99, 0xf9, 0x25, 0xe9, 0xa9, 0x89, 0x45, 0xfa, 0x60, 0x67, 0x24, 0xb1, 0x81, 0x29,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd4, 0xb1, 0xcb, 0xae, 0xb8, 0x00, 0x00, 0x00,
}

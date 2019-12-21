// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/app/proposal/delivery/grpc/response_grpc/response.proto

package response_grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Response struct {
	ID                   int64                `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	FreelancerId         int64                `protobuf:"varint,2,opt,name=FreelancerId,proto3" json:"FreelancerId,omitempty"`
	JobId                int64                `protobuf:"varint,3,opt,name=JobId,proto3" json:"JobId,omitempty"`
	Files                string               `protobuf:"bytes,4,opt,name=Files,proto3" json:"Files,omitempty"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,5,opt,name=Date,proto3" json:"Date,omitempty"`
	StatusManager        string               `protobuf:"bytes,6,opt,name=StatusManager,proto3" json:"StatusManager,omitempty"`
	StatusFreelancer     string               `protobuf:"bytes,7,opt,name=StatusFreelancer,proto3" json:"StatusFreelancer,omitempty"`
	PaymentAmount        float32              `protobuf:"fixed32,8,opt,name=PaymentAmount,proto3" json:"PaymentAmount,omitempty"`
	TimeEstimation       int32                `protobuf:"varint,9,opt,name=TimeEstimation,proto3" json:"TimeEstimation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ce2c2f34999c59e, []int{0}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Response) GetFreelancerId() int64 {
	if m != nil {
		return m.FreelancerId
	}
	return 0
}

func (m *Response) GetJobId() int64 {
	if m != nil {
		return m.JobId
	}
	return 0
}

func (m *Response) GetFiles() string {
	if m != nil {
		return m.Files
	}
	return ""
}

func (m *Response) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

func (m *Response) GetStatusManager() string {
	if m != nil {
		return m.StatusManager
	}
	return ""
}

func (m *Response) GetStatusFreelancer() string {
	if m != nil {
		return m.StatusFreelancer
	}
	return ""
}

func (m *Response) GetPaymentAmount() float32 {
	if m != nil {
		return m.PaymentAmount
	}
	return 0
}

func (m *Response) GetTimeEstimation() int32 {
	if m != nil {
		return m.TimeEstimation
	}
	return 0
}

type ResponseID struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseID) Reset()         { *m = ResponseID{} }
func (m *ResponseID) String() string { return proto.CompactTextString(m) }
func (*ResponseID) ProtoMessage()    {}
func (*ResponseID) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ce2c2f34999c59e, []int{1}
}

func (m *ResponseID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseID.Unmarshal(m, b)
}
func (m *ResponseID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseID.Marshal(b, m, deterministic)
}
func (m *ResponseID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseID.Merge(m, src)
}
func (m *ResponseID) XXX_Size() int {
	return xxx_messageInfo_ResponseID.Size(m)
}
func (m *ResponseID) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseID.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseID proto.InternalMessageInfo

func (m *ResponseID) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

type Status struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ce2c2f34999c59e, []int{2}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Nothing struct {
	Dummy                bool     `protobuf:"varint,1,opt,name=dummy,proto3" json:"dummy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Nothing) Reset()         { *m = Nothing{} }
func (m *Nothing) String() string { return proto.CompactTextString(m) }
func (*Nothing) ProtoMessage()    {}
func (*Nothing) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ce2c2f34999c59e, []int{3}
}

func (m *Nothing) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nothing.Unmarshal(m, b)
}
func (m *Nothing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nothing.Marshal(b, m, deterministic)
}
func (m *Nothing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nothing.Merge(m, src)
}
func (m *Nothing) XXX_Size() int {
	return xxx_messageInfo_Nothing.Size(m)
}
func (m *Nothing) XXX_DiscardUnknown() {
	xxx_messageInfo_Nothing.DiscardUnknown(m)
}

var xxx_messageInfo_Nothing proto.InternalMessageInfo

func (m *Nothing) GetDummy() bool {
	if m != nil {
		return m.Dummy
	}
	return false
}

func init() {
	proto.RegisterType((*Response)(nil), "response_grpc.Response")
	proto.RegisterType((*ResponseID)(nil), "response_grpc.ResponseID")
	proto.RegisterType((*Status)(nil), "response_grpc.Status")
	proto.RegisterType((*Nothing)(nil), "response_grpc.Nothing")
}

func init() {
	proto.RegisterFile("internal/app/proposal/delivery/grpc/response_grpc/response.proto", fileDescriptor_7ce2c2f34999c59e)
}

var fileDescriptor_7ce2c2f34999c59e = []byte{
	// 383 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x5f, 0x6b, 0xd4, 0x40,
	0x14, 0xc5, 0x49, 0x9a, 0x4d, 0x77, 0xaf, 0xb6, 0xca, 0x20, 0x3a, 0x86, 0x42, 0x43, 0x10, 0x09,
	0x3e, 0x24, 0x50, 0x9f, 0xf4, 0x49, 0x21, 0x2e, 0x46, 0x50, 0x64, 0xd4, 0x67, 0x99, 0x6d, 0xae,
	0x71, 0x20, 0xf3, 0x87, 0xc9, 0x44, 0xd8, 0x6f, 0xe0, 0x67, 0xf2, 0xd3, 0x49, 0x26, 0x0d, 0x4b,
	0xb6, 0xec, 0x5b, 0xee, 0xf9, 0x9d, 0xdc, 0x39, 0x9c, 0x0b, 0xef, 0x84, 0x72, 0x68, 0x15, 0xef,
	0x4a, 0x6e, 0x4c, 0x69, 0xac, 0x36, 0xba, 0xe7, 0x5d, 0xd9, 0x60, 0x27, 0xfe, 0xa0, 0xdd, 0x97,
	0xad, 0x35, 0xb7, 0xa5, 0xc5, 0xde, 0x68, 0xd5, 0xe3, 0xcf, 0xc5, 0x54, 0x18, 0xab, 0x9d, 0x26,
	0x17, 0x0b, 0x9a, 0x5c, 0xb7, 0x5a, 0xb7, 0x1d, 0x96, 0x1e, 0xee, 0x86, 0x5f, 0xa5, 0x13, 0x12,
	0x7b, 0xc7, 0xa5, 0x99, 0xfc, 0xd9, 0xbf, 0x10, 0xd6, 0xec, 0xee, 0x17, 0x72, 0x09, 0x61, 0x5d,
	0xd1, 0x20, 0x0d, 0xf2, 0x33, 0x16, 0xd6, 0x15, 0xc9, 0xe0, 0xe1, 0xd6, 0x22, 0x76, 0x5c, 0xdd,
	0xa2, 0xad, 0x1b, 0x1a, 0x7a, 0xb2, 0xd0, 0xc8, 0x13, 0x58, 0x7d, 0xd2, 0xbb, 0xba, 0xa1, 0x67,
	0x1e, 0x4e, 0xc3, 0xa8, 0x6e, 0x45, 0x87, 0x3d, 0x8d, 0xd2, 0x20, 0xdf, 0xb0, 0x69, 0x20, 0x05,
	0x44, 0x15, 0x77, 0x48, 0x57, 0x69, 0x90, 0x3f, 0xb8, 0x49, 0x8a, 0x29, 0x5c, 0x31, 0x87, 0x2b,
	0xbe, 0xcf, 0xe1, 0x98, 0xf7, 0x91, 0x17, 0x70, 0xf1, 0xcd, 0x71, 0x37, 0xf4, 0x9f, 0xb9, 0xe2,
	0x2d, 0x5a, 0x1a, 0xfb, 0x6d, 0x4b, 0x91, 0xbc, 0x82, 0xc7, 0x93, 0x70, 0xc8, 0x45, 0xcf, 0xbd,
	0xf1, 0x9e, 0x3e, 0x6e, 0xfc, 0xca, 0xf7, 0x12, 0x95, 0x7b, 0x2f, 0xf5, 0xa0, 0x1c, 0x5d, 0xa7,
	0x41, 0x1e, 0xb2, 0xa5, 0x48, 0x5e, 0xc2, 0xe5, 0x18, 0xe5, 0x43, 0xef, 0x84, 0xe4, 0x4e, 0x68,
	0x45, 0x37, 0x69, 0x90, 0xaf, 0xd8, 0x91, 0x9a, 0x5d, 0x01, 0xcc, 0xdd, 0xd5, 0xd5, 0x71, 0x7b,
	0xd9, 0x15, 0xc4, 0xd3, 0xfb, 0x84, 0x40, 0xa4, 0xb8, 0x44, 0xcf, 0x36, 0xcc, 0x7f, 0x67, 0xd7,
	0x70, 0xfe, 0x45, 0xbb, 0xdf, 0x42, 0xb5, 0x63, 0x59, 0xcd, 0x20, 0xe5, 0xde, 0xf3, 0x35, 0x9b,
	0x86, 0x9b, 0xbf, 0x01, 0x3c, 0x9a, 0xb7, 0x7f, 0xe4, 0xaa, 0xe9, 0xd0, 0x92, 0xb7, 0x10, 0x6d,
	0x85, 0x6a, 0xc8, 0xf3, 0x62, 0x71, 0xe6, 0xe2, 0x90, 0x22, 0x79, 0x76, 0x02, 0x91, 0x37, 0x10,
	0xff, 0x30, 0xcd, 0x58, 0xeb, 0x29, 0x4b, 0xf2, 0xf4, 0x08, 0xdc, 0x05, 0xdc, 0xc5, 0xfe, 0x42,
	0xaf, 0xff, 0x07, 0x00, 0x00, 0xff, 0xff, 0xd1, 0xa8, 0x44, 0x43, 0x9f, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ResponseHandlerClient is the client API for ResponseHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ResponseHandlerClient interface {
	Find(ctx context.Context, in *ResponseID, opts ...grpc.CallOption) (*Response, error)
	Update(ctx context.Context, in *Response, opts ...grpc.CallOption) (*Nothing, error)
}

type responseHandlerClient struct {
	cc *grpc.ClientConn
}

func NewResponseHandlerClient(cc *grpc.ClientConn) ResponseHandlerClient {
	return &responseHandlerClient{cc}
}

func (c *responseHandlerClient) Find(ctx context.Context, in *ResponseID, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/response_grpc.ResponseHandler/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *responseHandlerClient) Update(ctx context.Context, in *Response, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/response_grpc.ResponseHandler/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResponseHandlerServer is the server API for ResponseHandler service.
type ResponseHandlerServer interface {
	Find(context.Context, *ResponseID) (*Response, error)
	Update(context.Context, *Response) (*Nothing, error)
}

// UnimplementedResponseHandlerServer can be embedded to have forward compatible implementations.
type UnimplementedResponseHandlerServer struct {
}

func (*UnimplementedResponseHandlerServer) Find(ctx context.Context, req *ResponseID) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (*UnimplementedResponseHandlerServer) Update(ctx context.Context, req *Response) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func RegisterResponseHandlerServer(s *grpc.Server, srv ResponseHandlerServer) {
	s.RegisterService(&_ResponseHandler_serviceDesc, srv)
}

func _ResponseHandler_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResponseHandlerServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/response_grpc.ResponseHandler/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResponseHandlerServer).Find(ctx, req.(*ResponseID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResponseHandler_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Response)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResponseHandlerServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/response_grpc.ResponseHandler/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResponseHandlerServer).Update(ctx, req.(*Response))
	}
	return interceptor(ctx, in, info, handler)
}

var _ResponseHandler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "response_grpc.ResponseHandler",
	HandlerType: (*ResponseHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Find",
			Handler:    _ResponseHandler_Find_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ResponseHandler_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/proposal/delivery/grpc/response_grpc/response.proto",
}
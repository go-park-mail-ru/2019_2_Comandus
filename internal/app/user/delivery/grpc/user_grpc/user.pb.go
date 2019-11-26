// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/app/user/delivery/grpc/user_grpc/user.proto

package user_grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type UserID struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserID) Reset()         { *m = UserID{} }
func (m *UserID) String() string { return proto.CompactTextString(m) }
func (*UserID) ProtoMessage()    {}
func (*UserID) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9ce0a4e02b3337a, []int{0}
}

func (m *UserID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserID.Unmarshal(m, b)
}
func (m *UserID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserID.Marshal(b, m, deterministic)
}
func (m *UserID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserID.Merge(m, src)
}
func (m *UserID) XXX_Size() int {
	return xxx_messageInfo_UserID.Size(m)
}
func (m *UserID) XXX_DiscardUnknown() {
	xxx_messageInfo_UserID.DiscardUnknown(m)
}

var xxx_messageInfo_UserID proto.InternalMessageInfo

func (m *UserID) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

type User struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	FirstName            string   `protobuf:"bytes,2,opt,name=FirstName,proto3" json:"FirstName,omitempty"`
	SecondName           string   `protobuf:"bytes,3,opt,name=SecondName,proto3" json:"SecondName,omitempty"`
	UserName             string   `protobuf:"bytes,4,opt,name=UserName,proto3" json:"UserName,omitempty"`
	Email                string   `protobuf:"bytes,5,opt,name=Email,proto3" json:"Email,omitempty"`
	Password             string   `protobuf:"bytes,6,opt,name=Password,proto3" json:"Password,omitempty"`
	EncryptPassword      string   `protobuf:"bytes,7,opt,name=EncryptPassword,proto3" json:"EncryptPassword,omitempty"`
	UserType             string   `protobuf:"bytes,8,opt,name=UserType,proto3" json:"UserType,omitempty"`
	FreelancerId         int64    `protobuf:"varint,9,opt,name=FreelancerId,proto3" json:"FreelancerId,omitempty"`
	HireManagerId        int64    `protobuf:"varint,10,opt,name=HireManagerId,proto3" json:"HireManagerId,omitempty"`
	CompanyId            int64    `protobuf:"varint,11,opt,name=CompanyId,proto3" json:"CompanyId,omitempty"`
	Avatar               []byte   `protobuf:"bytes,12,opt,name=Avatar,proto3" json:"Avatar,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9ce0a4e02b3337a, []int{1}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *User) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *User) GetSecondName() string {
	if m != nil {
		return m.SecondName
	}
	return ""
}

func (m *User) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User) GetEncryptPassword() string {
	if m != nil {
		return m.EncryptPassword
	}
	return ""
}

func (m *User) GetUserType() string {
	if m != nil {
		return m.UserType
	}
	return ""
}

func (m *User) GetFreelancerId() int64 {
	if m != nil {
		return m.FreelancerId
	}
	return 0
}

func (m *User) GetHireManagerId() int64 {
	if m != nil {
		return m.HireManagerId
	}
	return 0
}

func (m *User) GetCompanyId() int64 {
	if m != nil {
		return m.CompanyId
	}
	return 0
}

func (m *User) GetAvatar() []byte {
	if m != nil {
		return m.Avatar
	}
	return nil
}

func init() {
	proto.RegisterType((*UserID)(nil), "user_grpc.UserID")
	proto.RegisterType((*User)(nil), "user_grpc.User")
}

func init() {
	proto.RegisterFile("internal/app/user/delivery/grpc/user_grpc/user.proto", fileDescriptor_b9ce0a4e02b3337a)
}

var fileDescriptor_b9ce0a4e02b3337a = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x49, 0xda, 0xc6, 0x66, 0x5a, 0x2d, 0x2e, 0x22, 0x4b, 0x11, 0x29, 0xc5, 0x43, 0xf0,
	0xd0, 0x80, 0x7a, 0xf1, 0x28, 0xb6, 0xa5, 0x39, 0x28, 0x52, 0xf5, 0x2c, 0x6b, 0x76, 0x28, 0x81,
	0x74, 0xb3, 0x4c, 0x62, 0x25, 0x6f, 0xec, 0x63, 0x48, 0xa6, 0x36, 0xb5, 0xb9, 0xe5, 0xff, 0xbe,
	0x7f, 0x33, 0xc9, 0x0e, 0xdc, 0x25, 0xa6, 0x40, 0x32, 0x2a, 0x0d, 0x95, 0xb5, 0xe1, 0x57, 0x8e,
	0x14, 0x6a, 0x4c, 0x93, 0x0d, 0x52, 0x19, 0xae, 0xc8, 0xc6, 0x8c, 0x3e, 0xea, 0xa7, 0x89, 0xa5,
	0xac, 0xc8, 0x84, 0x5f, 0xd3, 0xb1, 0x04, 0xef, 0x3d, 0x47, 0x8a, 0xa6, 0xe2, 0x04, 0xdc, 0x68,
	0x2a, 0x9d, 0x91, 0x13, 0xb4, 0x96, 0x6e, 0x34, 0x1d, 0xff, 0xb8, 0xd0, 0xae, 0x54, 0x53, 0x88,
	0x0b, 0xf0, 0xe7, 0x09, 0xe5, 0xc5, 0xb3, 0x5a, 0xa3, 0x74, 0x47, 0x4e, 0xe0, 0x2f, 0xf7, 0x40,
	0x5c, 0x02, 0xbc, 0x62, 0x9c, 0x19, 0xcd, 0xba, 0xc5, 0xfa, 0x1f, 0x11, 0x43, 0xe8, 0x56, 0x6f,
	0x65, 0xdb, 0x66, 0x5b, 0x67, 0x71, 0x06, 0x9d, 0xd9, 0x5a, 0x25, 0xa9, 0xec, 0xb0, 0xd8, 0x86,
	0xea, 0xc4, 0x8b, 0xca, 0xf3, 0xef, 0x8c, 0xb4, 0xf4, 0xb6, 0x27, 0x76, 0x59, 0x04, 0x30, 0x98,
	0x99, 0x98, 0x4a, 0x5b, 0xd4, 0x95, 0x23, 0xae, 0x34, 0xf1, 0x6e, 0xee, 0x5b, 0x69, 0x51, 0x76,
	0xf7, 0x73, 0xab, 0x2c, 0xc6, 0xd0, 0x9f, 0x13, 0x62, 0xaa, 0x4c, 0x8c, 0x14, 0x69, 0xe9, 0xf3,
	0xbf, 0x1e, 0x30, 0x71, 0x05, 0xc7, 0x8b, 0x84, 0xf0, 0x49, 0x19, 0xb5, 0xe2, 0x12, 0x70, 0xe9,
	0x10, 0x56, 0x77, 0xf3, 0x98, 0xad, 0xad, 0x32, 0x65, 0xa4, 0x65, 0x8f, 0x1b, 0x7b, 0x20, 0xce,
	0xc1, 0x7b, 0xd8, 0xa8, 0x42, 0x91, 0xec, 0x8f, 0x9c, 0xa0, 0xbf, 0xfc, 0x4b, 0x37, 0xf7, 0xd0,
	0xab, 0xbe, 0x65, 0xa1, 0x8c, 0x4e, 0x91, 0xc4, 0x35, 0xb4, 0xe7, 0x89, 0xd1, 0xe2, 0x74, 0x52,
	0xef, 0x69, 0xb2, 0x5d, 0xd2, 0x70, 0xd0, 0x40, 0x9f, 0x1e, 0x6f, 0xf4, 0xf6, 0x37, 0x00, 0x00,
	0xff, 0xff, 0xb2, 0xb8, 0xf5, 0xe5, 0x09, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UserHandlerClient is the client API for UserHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserHandlerClient interface {
	Find(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*User, error)
}

type userHandlerClient struct {
	cc *grpc.ClientConn
}

func NewUserHandlerClient(cc *grpc.ClientConn) UserHandlerClient {
	return &userHandlerClient{cc}
}

func (c *userHandlerClient) Find(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/user_grpc.UserHandler/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserHandlerServer is the server API for UserHandler service.
type UserHandlerServer interface {
	Find(context.Context, *UserID) (*User, error)
}

// UnimplementedUserHandlerServer can be embedded to have forward compatible implementations.
type UnimplementedUserHandlerServer struct {
}

func (*UnimplementedUserHandlerServer) Find(ctx context.Context, req *UserID) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}

func RegisterUserHandlerServer(s *grpc.Server, srv UserHandlerServer) {
	s.RegisterService(&_UserHandler_serviceDesc, srv)
}

func _UserHandler_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_grpc.UserHandler/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).Find(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserHandler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user_grpc.UserHandler",
	HandlerType: (*UserHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Find",
			Handler:    _UserHandler_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/user/delivery/grpc/user_grpc/user.proto",
}

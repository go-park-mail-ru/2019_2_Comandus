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
	return fileDescriptor_b9ce0a4e02b3337a, []int{0}
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
	return fileDescriptor_b9ce0a4e02b3337a, []int{1}
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

type Users struct {
	Names                []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Users) Reset()         { *m = Users{} }
func (m *Users) String() string { return proto.CompactTextString(m) }
func (*Users) ProtoMessage()    {}
func (*Users) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9ce0a4e02b3337a, []int{2}
}

func (m *Users) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Users.Unmarshal(m, b)
}
func (m *Users) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Users.Marshal(b, m, deterministic)
}
func (m *Users) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Users.Merge(m, src)
}
func (m *Users) XXX_Size() int {
	return xxx_messageInfo_Users.Size(m)
}
func (m *Users) XXX_DiscardUnknown() {
	xxx_messageInfo_Users.DiscardUnknown(m)
}

var xxx_messageInfo_Users proto.InternalMessageInfo

func (m *Users) GetNames() []string {
	if m != nil {
		return m.Names
	}
	return nil
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
	return fileDescriptor_b9ce0a4e02b3337a, []int{3}
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
	proto.RegisterType((*Nothing)(nil), "user_grpc.Nothing")
	proto.RegisterType((*UserID)(nil), "user_grpc.Freelancer")
	proto.RegisterType((*Users)(nil), "user_grpc.Users")
	proto.RegisterType((*User)(nil), "user_grpc.User")
}

func init() {
	proto.RegisterFile("internal/app/user/delivery/grpc/user_grpc/user.proto", fileDescriptor_b9ce0a4e02b3337a)
}

var fileDescriptor_b9ce0a4e02b3337a = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xdd, 0x4e, 0xa3, 0x40,
	0x1c, 0xc5, 0x03, 0x6d, 0x29, 0xfc, 0xdb, 0xdd, 0xee, 0x4e, 0x36, 0x9b, 0x49, 0xe3, 0x07, 0x21,
	0x5e, 0x10, 0x2f, 0x8a, 0x51, 0x5f, 0xc0, 0xd8, 0xd6, 0x72, 0x61, 0x63, 0x50, 0xaf, 0xcd, 0xc8,
	0x4c, 0x2a, 0x11, 0x06, 0x32, 0x43, 0x6b, 0x78, 0x63, 0x1f, 0xc3, 0xcc, 0xd0, 0xd2, 0x96, 0x3b,
	0xce, 0xf9, 0x9d, 0xc3, 0x7c, 0xfd, 0xe1, 0x36, 0xe1, 0x25, 0x13, 0x9c, 0xa4, 0x01, 0x29, 0x8a,
	0x60, 0x2d, 0x99, 0x08, 0x28, 0x4b, 0x93, 0x0d, 0x13, 0x55, 0xb0, 0x12, 0x45, 0xac, 0xad, 0xb7,
	0xe6, 0x6b, 0x52, 0x88, 0xbc, 0xcc, 0x91, 0xd3, 0xb8, 0xde, 0x39, 0xf4, 0x97, 0x79, 0xf9, 0x91,
	0xf0, 0x15, 0xfa, 0x07, 0x3d, 0xba, 0xce, 0xb2, 0x0a, 0x1b, 0xae, 0xe1, 0xdb, 0x51, 0x2d, 0x3c,
	0x0c, 0xd6, 0xab, 0x64, 0x22, 0x9c, 0xa2, 0xdf, 0x60, 0x86, 0x53, 0x0d, 0x3b, 0x91, 0x19, 0x4e,
	0xbd, 0x53, 0xe8, 0x29, 0x22, 0x55, 0x91, 0x93, 0x8c, 0x49, 0x6c, 0xb8, 0x1d, 0xdf, 0x89, 0x6a,
	0xe1, 0x7d, 0x9b, 0xd0, 0x55, 0xbc, 0xdd, 0x43, 0x27, 0xe0, 0xcc, 0x13, 0x21, 0xcb, 0x25, 0xc9,
	0x18, 0x36, 0x5d, 0xc3, 0x77, 0xa2, 0xbd, 0x81, 0xce, 0x00, 0x9e, 0x59, 0x9c, 0x73, 0xaa, 0x71,
	0x47, 0xe3, 0x03, 0x07, 0x8d, 0xc1, 0x56, 0x7f, 0xd5, 0xb4, 0xab, 0x69, 0xa3, 0xd5, 0x46, 0x66,
	0x19, 0x49, 0x52, 0xdc, 0xd3, 0xa0, 0x16, 0xaa, 0xf1, 0x44, 0xa4, 0xfc, 0xca, 0x05, 0xc5, 0x56,
	0xdd, 0xd8, 0x69, 0xe4, 0xc3, 0x68, 0xc6, 0x63, 0x51, 0x15, 0x65, 0x13, 0xe9, 0xeb, 0x48, 0xdb,
	0xde, 0xad, 0xfb, 0x52, 0x15, 0x0c, 0xdb, 0xfb, 0x75, 0x95, 0x46, 0x1e, 0x0c, 0xe7, 0x82, 0xb1,
	0x94, 0xf0, 0x98, 0x89, 0x90, 0x62, 0x47, 0x9f, 0xf5, 0xc8, 0x43, 0x17, 0xf0, 0x6b, 0x91, 0x08,
	0xf6, 0x48, 0x38, 0x59, 0xe9, 0x10, 0xe8, 0xd0, 0xb1, 0xa9, 0xee, 0xe6, 0x3e, 0xcf, 0x0a, 0xc2,
	0xab, 0x90, 0xe2, 0x81, 0x4e, 0xec, 0x0d, 0xf4, 0x1f, 0xac, 0xbb, 0x0d, 0x29, 0x89, 0xc0, 0x43,
	0xd7, 0xf0, 0x87, 0xd1, 0x56, 0x5d, 0x7f, 0xc2, 0x40, 0xed, 0x65, 0x41, 0x38, 0x4d, 0x99, 0x40,
	0x97, 0xd0, 0x9d, 0x27, 0x9c, 0xa2, 0xbf, 0x93, 0xe6, 0x9d, 0x27, 0xf5, 0x1b, 0x8e, 0x47, 0x2d,
	0x0b, 0x5d, 0x81, 0xfd, 0xc0, 0xf4, 0xcd, 0x4b, 0x84, 0x0e, 0xe0, 0x76, 0x28, 0xc6, 0x7f, 0x5a,
	0x05, 0xf9, 0x6e, 0xe9, 0x19, 0xba, 0xf9, 0x09, 0x00, 0x00, 0xff, 0xff, 0x92, 0x39, 0x21, 0xeb,
	0x7b, 0x02, 0x00, 0x00,
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
	GetNames(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Users, error)
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

func (c *userHandlerClient) GetNames(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Users, error) {
	out := new(Users)
	err := c.cc.Invoke(ctx, "/user_grpc.UserHandler/GetNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserHandlerServer is the server API for UserHandler service.
type UserHandlerServer interface {
	Find(context.Context, *UserID) (*User, error)
	GetNames(context.Context, *Nothing) (*Users, error)
}

// UnimplementedUserHandlerServer can be embedded to have forward compatible implementations.
type UnimplementedUserHandlerServer struct {
}

func (*UnimplementedUserHandlerServer) Find(ctx context.Context, req *UserID) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (*UnimplementedUserHandlerServer) GetNames(ctx context.Context, req *Nothing) (*Users, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNames not implemented")
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

func _UserHandler_GetNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserHandlerServer).GetNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_grpc.UserHandler/GetNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserHandlerServer).GetNames(ctx, req.(*Nothing))
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
		{
			MethodName: "GetNames",
			Handler:    _UserHandler_GetNames_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/user/delivery/grpc/user_grpc/user.proto",
}

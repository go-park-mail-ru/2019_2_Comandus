// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/app/company/delivery/grpc/company_grpc/company.proto

package company_grpc

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

type Company struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	CompanyName          string   `protobuf:"bytes,2,opt,name=CompanyName,proto3" json:"CompanyName,omitempty"`
	Site                 string   `protobuf:"bytes,3,opt,name=Site,proto3" json:"Site,omitempty"`
	TagLine              string   `protobuf:"bytes,4,opt,name=TagLine,proto3" json:"TagLine,omitempty"`
	Description          string   `protobuf:"bytes,5,opt,name=Description,proto3" json:"Description,omitempty"`
	Country              string   `protobuf:"bytes,6,opt,name=Country,proto3" json:"Country,omitempty"`
	City                 string   `protobuf:"bytes,7,opt,name=City,proto3" json:"City,omitempty"`
	Address              string   `protobuf:"bytes,8,opt,name=Address,proto3" json:"Address,omitempty"`
	Phone                string   `protobuf:"bytes,9,opt,name=Phone,proto3" json:"Phone,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Company) Reset()         { *m = Company{} }
func (m *Company) String() string { return proto.CompactTextString(m) }
func (*Company) ProtoMessage()    {}
func (*Company) Descriptor() ([]byte, []int) {
	return fileDescriptor_5755ad00c494f094, []int{0}
}

func (m *Company) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Company.Unmarshal(m, b)
}
func (m *Company) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Company.Marshal(b, m, deterministic)
}
func (m *Company) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Company.Merge(m, src)
}
func (m *Company) XXX_Size() int {
	return xxx_messageInfo_Company.Size(m)
}
func (m *Company) XXX_DiscardUnknown() {
	xxx_messageInfo_Company.DiscardUnknown(m)
}

var xxx_messageInfo_Company proto.InternalMessageInfo

func (m *Company) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Company) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *Company) GetSite() string {
	if m != nil {
		return m.Site
	}
	return ""
}

func (m *Company) GetTagLine() string {
	if m != nil {
		return m.TagLine
	}
	return ""
}

func (m *Company) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Company) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *Company) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *Company) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Company) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
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
	return fileDescriptor_5755ad00c494f094, []int{1}
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

func init() {
	proto.RegisterType((*Company)(nil), "company_grpc.Company")
	proto.RegisterType((*UserID)(nil), "company_grpc.UserID")
}

func init() {
	proto.RegisterFile("internal/app/company/delivery/grpc/company_grpc/company.proto", fileDescriptor_5755ad00c494f094)
}

var fileDescriptor_5755ad00c494f094 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0x49, 0xff, 0x24, 0x76, 0xd4, 0x1e, 0x96, 0x0a, 0x83, 0xa7, 0xd0, 0x53, 0x4f, 0x09,
	0xe8, 0x55, 0x0f, 0x92, 0x1c, 0x0c, 0x48, 0x91, 0xaa, 0x67, 0x59, 0x93, 0xa1, 0x2e, 0xa4, 0xbb,
	0xcb, 0x64, 0x15, 0xf2, 0x91, 0xfd, 0x16, 0x92, 0x6c, 0x02, 0xd1, 0xdb, 0xbc, 0xdf, 0xcc, 0xbc,
	0x61, 0x1e, 0xdc, 0x2b, 0xed, 0x88, 0xb5, 0xac, 0x53, 0x69, 0x6d, 0x5a, 0x9a, 0x93, 0x95, 0xba,
	0x4d, 0x2b, 0xaa, 0xd5, 0x37, 0x71, 0x9b, 0x1e, 0xd9, 0x96, 0x23, 0x7d, 0x9f, 0x8a, 0xc4, 0xb2,
	0x71, 0x46, 0x5c, 0x4c, 0x7b, 0xdb, 0x9f, 0x00, 0xa2, 0xcc, 0x03, 0xb1, 0x86, 0x59, 0x91, 0x63,
	0x10, 0x07, 0xbb, 0xf9, 0x61, 0x56, 0xe4, 0x22, 0x86, 0xf3, 0xa1, 0xb5, 0x97, 0x27, 0xc2, 0x59,
	0x1c, 0xec, 0x56, 0x87, 0x29, 0x12, 0x02, 0x16, 0x2f, 0xca, 0x11, 0xce, 0xfb, 0x56, 0x5f, 0x0b,
	0x84, 0xe8, 0x55, 0x1e, 0x9f, 0x94, 0x26, 0x5c, 0xf4, 0x78, 0x94, 0x9d, 0x5f, 0x4e, 0x4d, 0xc9,
	0xca, 0x3a, 0x65, 0x34, 0x2e, 0xbd, 0xdf, 0x04, 0x75, 0xbb, 0x99, 0xf9, 0xd2, 0x8e, 0x5b, 0x0c,
	0xfd, 0xee, 0x20, 0xbb, 0x4b, 0x99, 0x72, 0x2d, 0x46, 0xfe, 0x52, 0x57, 0x77, 0xd3, 0x0f, 0x55,
	0xc5, 0xd4, 0x34, 0x78, 0xe6, 0xa7, 0x07, 0x29, 0x36, 0xb0, 0x7c, 0xfe, 0x34, 0x9a, 0x70, 0xd5,
	0x73, 0x2f, 0xb6, 0x08, 0xe1, 0x5b, 0x43, 0x5c, 0xe4, 0xff, 0x3f, 0xbd, 0xd9, 0xc3, 0x7a, 0x78,
	0xeb, 0x51, 0xea, 0xaa, 0x26, 0x16, 0x77, 0x70, 0x99, 0x31, 0x49, 0x47, 0x63, 0x38, 0x9b, 0x64,
	0x9a, 0x5b, 0xe2, 0x8d, 0xae, 0xaf, 0xfe, 0xd2, 0x61, 0xf8, 0x23, 0xec, 0xa3, 0xbe, 0xfd, 0x0d,
	0x00, 0x00, 0xff, 0xff, 0x03, 0x6c, 0x9b, 0xb1, 0xab, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CompanyHandlerClient is the client API for CompanyHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CompanyHandlerClient interface {
	CreateCompany(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Company, error)
}

type companyHandlerClient struct {
	cc *grpc.ClientConn
}

func NewCompanyHandlerClient(cc *grpc.ClientConn) CompanyHandlerClient {
	return &companyHandlerClient{cc}
}

func (c *companyHandlerClient) CreateCompany(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Company, error) {
	out := new(Company)
	err := c.cc.Invoke(ctx, "/company_grpc.CompanyHandler/CreateCompany", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompanyHandlerServer is the server API for CompanyHandler service.
type CompanyHandlerServer interface {
	CreateCompany(context.Context, *UserID) (*Company, error)
}

// UnimplementedCompanyHandlerServer can be embedded to have forward compatible implementations.
type UnimplementedCompanyHandlerServer struct {
}

func (*UnimplementedCompanyHandlerServer) CreateCompany(ctx context.Context, req *UserID) (*Company, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCompany not implemented")
}

func RegisterCompanyHandlerServer(s *grpc.Server, srv CompanyHandlerServer) {
	s.RegisterService(&_CompanyHandler_serviceDesc, srv)
}

func _CompanyHandler_CreateCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompanyHandlerServer).CreateCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company_grpc.CompanyHandler/CreateCompany",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompanyHandlerServer).CreateCompany(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

var _CompanyHandler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "company_grpc.CompanyHandler",
	HandlerType: (*CompanyHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCompany",
			Handler:    _CompanyHandler_CreateCompany_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/company/delivery/grpc/company_grpc/company.proto",
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/app/location/delivery/grpc/location_grpc/location.proto

package location_grpc

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

type CountryID struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CountryID) Reset()         { *m = CountryID{} }
func (m *CountryID) String() string { return proto.CompactTextString(m) }
func (*CountryID) ProtoMessage()    {}
func (*CountryID) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca09ea382dbdfc92, []int{0}
}

func (m *CountryID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CountryID.Unmarshal(m, b)
}
func (m *CountryID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CountryID.Marshal(b, m, deterministic)
}
func (m *CountryID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CountryID.Merge(m, src)
}
func (m *CountryID) XXX_Size() int {
	return xxx_messageInfo_CountryID.Size(m)
}
func (m *CountryID) XXX_DiscardUnknown() {
	xxx_messageInfo_CountryID.DiscardUnknown(m)
}

var xxx_messageInfo_CountryID proto.InternalMessageInfo

func (m *CountryID) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

type CityID struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CityID) Reset()         { *m = CityID{} }
func (m *CityID) String() string { return proto.CompactTextString(m) }
func (*CityID) ProtoMessage()    {}
func (*CityID) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca09ea382dbdfc92, []int{1}
}

func (m *CityID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CityID.Unmarshal(m, b)
}
func (m *CityID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CityID.Marshal(b, m, deterministic)
}
func (m *CityID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CityID.Merge(m, src)
}
func (m *CityID) XXX_Size() int {
	return xxx_messageInfo_CityID.Size(m)
}
func (m *CityID) XXX_DiscardUnknown() {
	xxx_messageInfo_CityID.DiscardUnknown(m)
}

var xxx_messageInfo_CityID proto.InternalMessageInfo

func (m *CityID) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

type Country struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Country) Reset()         { *m = Country{} }
func (m *Country) String() string { return proto.CompactTextString(m) }
func (*Country) ProtoMessage()    {}
func (*Country) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca09ea382dbdfc92, []int{2}
}

func (m *Country) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Country.Unmarshal(m, b)
}
func (m *Country) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Country.Marshal(b, m, deterministic)
}
func (m *Country) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Country.Merge(m, src)
}
func (m *Country) XXX_Size() int {
	return xxx_messageInfo_Country.Size(m)
}
func (m *Country) XXX_DiscardUnknown() {
	xxx_messageInfo_Country.DiscardUnknown(m)
}

var xxx_messageInfo_Country proto.InternalMessageInfo

func (m *Country) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Country) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type City struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	CountryID            int64    `protobuf:"varint,2,opt,name=CountryID,proto3" json:"CountryID,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *City) Reset()         { *m = City{} }
func (m *City) String() string { return proto.CompactTextString(m) }
func (*City) ProtoMessage()    {}
func (*City) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca09ea382dbdfc92, []int{3}
}

func (m *City) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_City.Unmarshal(m, b)
}
func (m *City) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_City.Marshal(b, m, deterministic)
}
func (m *City) XXX_Merge(src proto.Message) {
	xxx_messageInfo_City.Merge(m, src)
}
func (m *City) XXX_Size() int {
	return xxx_messageInfo_City.Size(m)
}
func (m *City) XXX_DiscardUnknown() {
	xxx_messageInfo_City.DiscardUnknown(m)
}

var xxx_messageInfo_City proto.InternalMessageInfo

func (m *City) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *City) GetCountryID() int64 {
	if m != nil {
		return m.CountryID
	}
	return 0
}

func (m *City) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*CountryID)(nil), "location_grpc.CountryID")
	proto.RegisterType((*CityID)(nil), "location_grpc.CityID")
	proto.RegisterType((*Country)(nil), "location_grpc.Country")
	proto.RegisterType((*City)(nil), "location_grpc.City")
}

func init() {
	proto.RegisterFile("internal/app/location/delivery/grpc/location_grpc/location.proto", fileDescriptor_ca09ea382dbdfc92)
}

var fileDescriptor_ca09ea382dbdfc92 = []byte{
	// 220 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0xc8, 0xcc, 0x2b, 0x49,
	0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x4f, 0x2c, 0x28, 0xd0, 0xcf, 0xc9, 0x4f, 0x4e, 0x2c, 0xc9, 0xcc,
	0xcf, 0xd3, 0x4f, 0x49, 0xcd, 0xc9, 0x2c, 0x4b, 0x2d, 0xaa, 0xd4, 0x4f, 0x2f, 0x2a, 0x48, 0x86,
	0x0b, 0xc7, 0xa3, 0xf0, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x78, 0x51, 0x64, 0x95, 0xa4,
	0xb9, 0x38, 0x9d, 0xf3, 0x4b, 0xf3, 0x4a, 0x8a, 0x2a, 0x3d, 0x5d, 0x84, 0xf8, 0xb8, 0x98, 0x3c,
	0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x98, 0x83, 0x98, 0x3c, 0x5d, 0x94, 0x24, 0xb8, 0xd8, 0x9c,
	0x33, 0x4b, 0xb0, 0xc9, 0xe8, 0x72, 0xb1, 0x43, 0xb5, 0xa1, 0x4b, 0x09, 0x09, 0x71, 0xb1, 0xf8,
	0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0x1e, 0x5c, 0x2c,
	0x20, 0x83, 0x30, 0xd4, 0xca, 0x20, 0xd9, 0x0e, 0xd6, 0xc0, 0x1c, 0x84, 0xe4, 0x1c, 0x98, 0x49,
	0xcc, 0x08, 0x93, 0x8c, 0x3a, 0x18, 0xb9, 0xf8, 0x7d, 0xa0, 0x3e, 0xf0, 0x48, 0xcc, 0x4b, 0xc9,
	0x49, 0x2d, 0x12, 0xb2, 0xe3, 0xe2, 0x72, 0x4f, 0x2d, 0x81, 0xb9, 0x47, 0x42, 0x0f, 0xc5, 0x87,
	0x7a, 0x70, 0xf3, 0xa4, 0xc4, 0xb0, 0xcb, 0x08, 0x99, 0x72, 0xb1, 0x83, 0xf4, 0x83, 0x1c, 0x28,
	0x8a, 0xae, 0x04, 0xec, 0x7d, 0x29, 0x61, 0x2c, 0xc2, 0x49, 0x6c, 0xe0, 0x00, 0x35, 0x06, 0x04,
	0x00, 0x00, 0xff, 0xff, 0x7f, 0x2e, 0x4e, 0x1c, 0x94, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LocationHandlerClient is the client API for LocationHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LocationHandlerClient interface {
	GetCountry(ctx context.Context, in *CountryID, opts ...grpc.CallOption) (*Country, error)
	GetCity(ctx context.Context, in *CityID, opts ...grpc.CallOption) (*City, error)
}

type locationHandlerClient struct {
	cc *grpc.ClientConn
}

func NewLocationHandlerClient(cc *grpc.ClientConn) LocationHandlerClient {
	return &locationHandlerClient{cc}
}

func (c *locationHandlerClient) GetCountry(ctx context.Context, in *CountryID, opts ...grpc.CallOption) (*Country, error) {
	out := new(Country)
	err := c.cc.Invoke(ctx, "/location_grpc.LocationHandler/GetCountry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationHandlerClient) GetCity(ctx context.Context, in *CityID, opts ...grpc.CallOption) (*City, error) {
	out := new(City)
	err := c.cc.Invoke(ctx, "/location_grpc.LocationHandler/GetCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocationHandlerServer is the server API for LocationHandler service.
type LocationHandlerServer interface {
	GetCountry(context.Context, *CountryID) (*Country, error)
	GetCity(context.Context, *CityID) (*City, error)
}

// UnimplementedLocationHandlerServer can be embedded to have forward compatible implementations.
type UnimplementedLocationHandlerServer struct {
}

func (*UnimplementedLocationHandlerServer) GetCountry(ctx context.Context, req *CountryID) (*Country, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCountry not implemented")
}
func (*UnimplementedLocationHandlerServer) GetCity(ctx context.Context, req *CityID) (*City, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCity not implemented")
}

func RegisterLocationHandlerServer(s *grpc.Server, srv LocationHandlerServer) {
	s.RegisterService(&_LocationHandler_serviceDesc, srv)
}

func _LocationHandler_GetCountry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountryID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationHandlerServer).GetCountry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location_grpc.LocationHandler/GetCountry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationHandlerServer).GetCountry(ctx, req.(*CountryID))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocationHandler_GetCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CityID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationHandlerServer).GetCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location_grpc.LocationHandler/GetCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationHandlerServer).GetCity(ctx, req.(*CityID))
	}
	return interceptor(ctx, in, info, handler)
}

var _LocationHandler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "location_grpc.LocationHandler",
	HandlerType: (*LocationHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCountry",
			Handler:    _LocationHandler_GetCountry_Handler,
		},
		{
			MethodName: "GetCity",
			Handler:    _LocationHandler_GetCity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/location/delivery/grpc/location_grpc/location.proto",
}
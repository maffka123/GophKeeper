// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: server.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GophKeeperClient is the client API for GophKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResp, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResp, error)
	Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*InsertResp, error)
	GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*GetDataResp, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResp, error)
	GetAllDataForUser(ctx context.Context, in *GetAllDataForUserRequest, opts ...grpc.CallOption) (*GetAllDataForUserResp, error)
	InsertSyncData(ctx context.Context, in *InsertSyncDataRequest, opts ...grpc.CallOption) (*InsertSyncDataResp, error)
}

type gophKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperClient(cc grpc.ClientConnInterface) GophKeeperClient {
	return &gophKeeperClient{cc}
}

func (c *gophKeeperClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResp, error) {
	out := new(RegisterResp)
	err := c.cc.Invoke(ctx, "/proto.GophKeeper/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResp, error) {
	out := new(LoginResp)
	err := c.cc.Invoke(ctx, "/proto.GophKeeper/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*InsertResp, error) {
	out := new(InsertResp)
	err := c.cc.Invoke(ctx, "/proto.GophKeeper/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*GetDataResp, error) {
	out := new(GetDataResp)
	err := c.cc.Invoke(ctx, "/proto.GophKeeper/GetData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResp, error) {
	out := new(DeleteResp)
	err := c.cc.Invoke(ctx, "/proto.GophKeeper/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) GetAllDataForUser(ctx context.Context, in *GetAllDataForUserRequest, opts ...grpc.CallOption) (*GetAllDataForUserResp, error) {
	out := new(GetAllDataForUserResp)
	err := c.cc.Invoke(ctx, "/proto.GophKeeper/GetAllDataForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) InsertSyncData(ctx context.Context, in *InsertSyncDataRequest, opts ...grpc.CallOption) (*InsertSyncDataResp, error) {
	out := new(InsertSyncDataResp)
	err := c.cc.Invoke(ctx, "/proto.GophKeeper/InsertSyncData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServer is the server API for GophKeeper service.
// All implementations must embed UnimplementedGophKeeperServer
// for forward compatibility
type GophKeeperServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResp, error)
	Login(context.Context, *LoginRequest) (*LoginResp, error)
	Insert(context.Context, *InsertRequest) (*InsertResp, error)
	GetData(context.Context, *GetDataRequest) (*GetDataResp, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResp, error)
	GetAllDataForUser(context.Context, *GetAllDataForUserRequest) (*GetAllDataForUserResp, error)
	InsertSyncData(context.Context, *InsertSyncDataRequest) (*InsertSyncDataResp, error)
	mustEmbedUnimplementedGophKeeperServer()
}

// UnimplementedGophKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedGophKeeperServer struct {
}

func (UnimplementedGophKeeperServer) Register(context.Context, *RegisterRequest) (*RegisterResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedGophKeeperServer) Login(context.Context, *LoginRequest) (*LoginResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedGophKeeperServer) Insert(context.Context, *InsertRequest) (*InsertResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (UnimplementedGophKeeperServer) GetData(context.Context, *GetDataRequest) (*GetDataResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedGophKeeperServer) Delete(context.Context, *DeleteRequest) (*DeleteResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedGophKeeperServer) GetAllDataForUser(context.Context, *GetAllDataForUserRequest) (*GetAllDataForUserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllDataForUser not implemented")
}
func (UnimplementedGophKeeperServer) InsertSyncData(context.Context, *InsertSyncDataRequest) (*InsertSyncDataResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertSyncData not implemented")
}
func (UnimplementedGophKeeperServer) mustEmbedUnimplementedGophKeeperServer() {}

// UnsafeGophKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServer will
// result in compilation errors.
type UnsafeGophKeeperServer interface {
	mustEmbedUnimplementedGophKeeperServer()
}

func RegisterGophKeeperServer(s grpc.ServiceRegistrar, srv GophKeeperServer) {
	s.RegisterService(&GophKeeper_ServiceDesc, srv)
}

func _GophKeeper_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeper/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeper/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeper/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Insert(ctx, req.(*InsertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeper/GetData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).GetData(ctx, req.(*GetDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeper/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_GetAllDataForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllDataForUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).GetAllDataForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeper/GetAllDataForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).GetAllDataForUser(ctx, req.(*GetAllDataForUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_InsertSyncData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertSyncDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).InsertSyncData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GophKeeper/InsertSyncData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).InsertSyncData(ctx, req.(*InsertSyncDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeper_ServiceDesc is the grpc.ServiceDesc for GophKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GophKeeper",
	HandlerType: (*GophKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _GophKeeper_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _GophKeeper_Login_Handler,
		},
		{
			MethodName: "Insert",
			Handler:    _GophKeeper_Insert_Handler,
		},
		{
			MethodName: "GetData",
			Handler:    _GophKeeper_GetData_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _GophKeeper_Delete_Handler,
		},
		{
			MethodName: "GetAllDataForUser",
			Handler:    _GophKeeper_GetAllDataForUser_Handler,
		},
		{
			MethodName: "InsertSyncData",
			Handler:    _GophKeeper_InsertSyncData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: proto/gophkeeper.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Gophkeeper_RegisterUser_FullMethodName       = "/gophkeeper.Gophkeeper/RegisterUser"
	Gophkeeper_Authorize_FullMethodName          = "/gophkeeper.Gophkeeper/Authorize"
	Gophkeeper_Echo_FullMethodName               = "/gophkeeper.Gophkeeper/Echo"
	Gophkeeper_AddBankCard_FullMethodName        = "/gophkeeper.Gophkeeper/AddBankCard"
	Gophkeeper_RemoveBankCard_FullMethodName     = "/gophkeeper.Gophkeeper/RemoveBankCard"
	Gophkeeper_GetBankCards_FullMethodName       = "/gophkeeper.Gophkeeper/GetBankCards"
	Gophkeeper_GetBankCard_FullMethodName        = "/gophkeeper.Gophkeeper/GetBankCard"
	Gophkeeper_AddUserCredentials_FullMethodName = "/gophkeeper.Gophkeeper/AddUserCredentials"
	Gophkeeper_GetUserCredentials_FullMethodName = "/gophkeeper.Gophkeeper/GetUserCredentials"
	Gophkeeper_GetUserCredential_FullMethodName  = "/gophkeeper.Gophkeeper/GetUserCredential"
	Gophkeeper_Upload_FullMethodName             = "/gophkeeper.Gophkeeper/Upload"
	Gophkeeper_Download_FullMethodName           = "/gophkeeper.Gophkeeper/Download"
	Gophkeeper_GetFiles_FullMethodName           = "/gophkeeper.Gophkeeper/GetFiles"
)

// GophkeeperClient is the client API for Gophkeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophkeeperClient interface {
	RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*AuthorizeResponse, error)
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	AddBankCard(ctx context.Context, in *AddBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RemoveBankCard(ctx context.Context, in *RemoveBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetBankCards(ctx context.Context, in *GetBankCardsRequest, opts ...grpc.CallOption) (*GetBankCardsResponse, error)
	GetBankCard(ctx context.Context, in *GetBankCardRequest, opts ...grpc.CallOption) (*GetBankCardResponse, error)
	AddUserCredentials(ctx context.Context, in *AddUserCredentialsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetUserCredentials(ctx context.Context, in *GetUserCredentialsRequest, opts ...grpc.CallOption) (*GetUserCredentialsResponse, error)
	GetUserCredential(ctx context.Context, in *GetUserCredentialRequest, opts ...grpc.CallOption) (*GetUserCredentialResponse, error)
	Upload(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileUploadRequest, FileUploadResponse], error)
	Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[FileDownloadResponse], error)
	GetFiles(ctx context.Context, in *GetFilesRequest, opts ...grpc.CallOption) (*GetFilesResponse, error)
}

type gophkeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophkeeperClient(cc grpc.ClientConnInterface) GophkeeperClient {
	return &gophkeeperClient{cc}
}

func (c *gophkeeperClient) RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gophkeeper_RegisterUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*AuthorizeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthorizeResponse)
	err := c.cc.Invoke(ctx, Gophkeeper_Authorize_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, Gophkeeper_Echo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) AddBankCard(ctx context.Context, in *AddBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gophkeeper_AddBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) RemoveBankCard(ctx context.Context, in *RemoveBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gophkeeper_RemoveBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetBankCards(ctx context.Context, in *GetBankCardsRequest, opts ...grpc.CallOption) (*GetBankCardsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBankCardsResponse)
	err := c.cc.Invoke(ctx, Gophkeeper_GetBankCards_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetBankCard(ctx context.Context, in *GetBankCardRequest, opts ...grpc.CallOption) (*GetBankCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBankCardResponse)
	err := c.cc.Invoke(ctx, Gophkeeper_GetBankCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) AddUserCredentials(ctx context.Context, in *AddUserCredentialsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gophkeeper_AddUserCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetUserCredentials(ctx context.Context, in *GetUserCredentialsRequest, opts ...grpc.CallOption) (*GetUserCredentialsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserCredentialsResponse)
	err := c.cc.Invoke(ctx, Gophkeeper_GetUserCredentials_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetUserCredential(ctx context.Context, in *GetUserCredentialRequest, opts ...grpc.CallOption) (*GetUserCredentialResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserCredentialResponse)
	err := c.cc.Invoke(ctx, Gophkeeper_GetUserCredential_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) Upload(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileUploadRequest, FileUploadResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Gophkeeper_ServiceDesc.Streams[0], Gophkeeper_Upload_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileUploadRequest, FileUploadResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Gophkeeper_UploadClient = grpc.ClientStreamingClient[FileUploadRequest, FileUploadResponse]

func (c *gophkeeperClient) Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[FileDownloadResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Gophkeeper_ServiceDesc.Streams[1], Gophkeeper_Download_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileDownloadRequest, FileDownloadResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Gophkeeper_DownloadClient = grpc.ServerStreamingClient[FileDownloadResponse]

func (c *gophkeeperClient) GetFiles(ctx context.Context, in *GetFilesRequest, opts ...grpc.CallOption) (*GetFilesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFilesResponse)
	err := c.cc.Invoke(ctx, Gophkeeper_GetFiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophkeeperServer is the server API for Gophkeeper service.
// All implementations must embed UnimplementedGophkeeperServer
// for forward compatibility.
type GophkeeperServer interface {
	RegisterUser(context.Context, *RegisterUserRequest) (*emptypb.Empty, error)
	Authorize(context.Context, *AuthorizeRequest) (*AuthorizeResponse, error)
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	AddBankCard(context.Context, *AddBankCardRequest) (*emptypb.Empty, error)
	RemoveBankCard(context.Context, *RemoveBankCardRequest) (*emptypb.Empty, error)
	GetBankCards(context.Context, *GetBankCardsRequest) (*GetBankCardsResponse, error)
	GetBankCard(context.Context, *GetBankCardRequest) (*GetBankCardResponse, error)
	AddUserCredentials(context.Context, *AddUserCredentialsRequest) (*emptypb.Empty, error)
	GetUserCredentials(context.Context, *GetUserCredentialsRequest) (*GetUserCredentialsResponse, error)
	GetUserCredential(context.Context, *GetUserCredentialRequest) (*GetUserCredentialResponse, error)
	Upload(grpc.ClientStreamingServer[FileUploadRequest, FileUploadResponse]) error
	Download(*FileDownloadRequest, grpc.ServerStreamingServer[FileDownloadResponse]) error
	GetFiles(context.Context, *GetFilesRequest) (*GetFilesResponse, error)
	mustEmbedUnimplementedGophkeeperServer()
}

// UnimplementedGophkeeperServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGophkeeperServer struct{}

func (UnimplementedGophkeeperServer) RegisterUser(context.Context, *RegisterUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedGophkeeperServer) Authorize(context.Context, *AuthorizeRequest) (*AuthorizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (UnimplementedGophkeeperServer) Echo(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedGophkeeperServer) AddBankCard(context.Context, *AddBankCardRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBankCard not implemented")
}
func (UnimplementedGophkeeperServer) RemoveBankCard(context.Context, *RemoveBankCardRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveBankCard not implemented")
}
func (UnimplementedGophkeeperServer) GetBankCards(context.Context, *GetBankCardsRequest) (*GetBankCardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBankCards not implemented")
}
func (UnimplementedGophkeeperServer) GetBankCard(context.Context, *GetBankCardRequest) (*GetBankCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBankCard not implemented")
}
func (UnimplementedGophkeeperServer) AddUserCredentials(context.Context, *AddUserCredentialsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserCredentials not implemented")
}
func (UnimplementedGophkeeperServer) GetUserCredentials(context.Context, *GetUserCredentialsRequest) (*GetUserCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserCredentials not implemented")
}
func (UnimplementedGophkeeperServer) GetUserCredential(context.Context, *GetUserCredentialRequest) (*GetUserCredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserCredential not implemented")
}
func (UnimplementedGophkeeperServer) Upload(grpc.ClientStreamingServer[FileUploadRequest, FileUploadResponse]) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedGophkeeperServer) Download(*FileDownloadRequest, grpc.ServerStreamingServer[FileDownloadResponse]) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedGophkeeperServer) GetFiles(context.Context, *GetFilesRequest) (*GetFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFiles not implemented")
}
func (UnimplementedGophkeeperServer) mustEmbedUnimplementedGophkeeperServer() {}
func (UnimplementedGophkeeperServer) testEmbeddedByValue()                    {}

// UnsafeGophkeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophkeeperServer will
// result in compilation errors.
type UnsafeGophkeeperServer interface {
	mustEmbedUnimplementedGophkeeperServer()
}

func RegisterGophkeeperServer(s grpc.ServiceRegistrar, srv GophkeeperServer) {
	// If the following call pancis, it indicates UnimplementedGophkeeperServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Gophkeeper_ServiceDesc, srv)
}

func _Gophkeeper_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_RegisterUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).RegisterUser(ctx, req.(*RegisterUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_Authorize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).Authorize(ctx, req.(*AuthorizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_Echo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_AddBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).AddBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_AddBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).AddBankCard(ctx, req.(*AddBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_RemoveBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).RemoveBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_RemoveBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).RemoveBankCard(ctx, req.(*RemoveBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetBankCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBankCardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetBankCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_GetBankCards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetBankCards(ctx, req.(*GetBankCardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_GetBankCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetBankCard(ctx, req.(*GetBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_AddUserCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).AddUserCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_AddUserCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).AddUserCredentials(ctx, req.(*AddUserCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetUserCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetUserCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_GetUserCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetUserCredentials(ctx, req.(*GetUserCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetUserCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserCredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetUserCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_GetUserCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetUserCredential(ctx, req.(*GetUserCredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GophkeeperServer).Upload(&grpc.GenericServerStream[FileUploadRequest, FileUploadResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Gophkeeper_UploadServer = grpc.ClientStreamingServer[FileUploadRequest, FileUploadResponse]

func _Gophkeeper_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileDownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GophkeeperServer).Download(m, &grpc.GenericServerStream[FileDownloadRequest, FileDownloadResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Gophkeeper_DownloadServer = grpc.ServerStreamingServer[FileDownloadResponse]

func _Gophkeeper_GetFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gophkeeper_GetFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetFiles(ctx, req.(*GetFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gophkeeper_ServiceDesc is the grpc.ServiceDesc for Gophkeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gophkeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gophkeeper.Gophkeeper",
	HandlerType: (*GophkeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterUser",
			Handler:    _Gophkeeper_RegisterUser_Handler,
		},
		{
			MethodName: "Authorize",
			Handler:    _Gophkeeper_Authorize_Handler,
		},
		{
			MethodName: "Echo",
			Handler:    _Gophkeeper_Echo_Handler,
		},
		{
			MethodName: "AddBankCard",
			Handler:    _Gophkeeper_AddBankCard_Handler,
		},
		{
			MethodName: "RemoveBankCard",
			Handler:    _Gophkeeper_RemoveBankCard_Handler,
		},
		{
			MethodName: "GetBankCards",
			Handler:    _Gophkeeper_GetBankCards_Handler,
		},
		{
			MethodName: "GetBankCard",
			Handler:    _Gophkeeper_GetBankCard_Handler,
		},
		{
			MethodName: "AddUserCredentials",
			Handler:    _Gophkeeper_AddUserCredentials_Handler,
		},
		{
			MethodName: "GetUserCredentials",
			Handler:    _Gophkeeper_GetUserCredentials_Handler,
		},
		{
			MethodName: "GetUserCredential",
			Handler:    _Gophkeeper_GetUserCredential_Handler,
		},
		{
			MethodName: "GetFiles",
			Handler:    _Gophkeeper_GetFiles_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _Gophkeeper_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _Gophkeeper_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/gophkeeper.proto",
}

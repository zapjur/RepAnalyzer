// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: proto/db.proto

package db

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DBService_GetUser_FullMethodName                 = "/db.DBService/GetUser"
	DBService_SaveUploadedVideo_FullMethodName       = "/db.DBService/SaveUploadedVideo"
	DBService_GetUserVideosByExercise_FullMethodName = "/db.DBService/GetUserVideosByExercise"
	DBService_SaveAnalysis_FullMethodName            = "/db.DBService/SaveAnalysis"
)

// DBServiceClient is the client API for DBService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DBServiceClient interface {
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	SaveUploadedVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error)
	GetUserVideosByExercise(ctx context.Context, in *GetUserVideosByExerciseRequest, opts ...grpc.CallOption) (*GetUserVideosByExerciseResponse, error)
	SaveAnalysis(ctx context.Context, in *VideoAnalysisRequest, opts ...grpc.CallOption) (*SaveAnalysisResponse, error)
}

type dBServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDBServiceClient(cc grpc.ClientConnInterface) DBServiceClient {
	return &dBServiceClient{cc}
}

func (c *dBServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, DBService_GetUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) SaveUploadedVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadVideoResponse)
	err := c.cc.Invoke(ctx, DBService_SaveUploadedVideo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) GetUserVideosByExercise(ctx context.Context, in *GetUserVideosByExerciseRequest, opts ...grpc.CallOption) (*GetUserVideosByExerciseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserVideosByExerciseResponse)
	err := c.cc.Invoke(ctx, DBService_GetUserVideosByExercise_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) SaveAnalysis(ctx context.Context, in *VideoAnalysisRequest, opts ...grpc.CallOption) (*SaveAnalysisResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SaveAnalysisResponse)
	err := c.cc.Invoke(ctx, DBService_SaveAnalysis_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DBServiceServer is the server API for DBService service.
// All implementations must embed UnimplementedDBServiceServer
// for forward compatibility.
type DBServiceServer interface {
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	SaveUploadedVideo(context.Context, *UploadVideoRequest) (*UploadVideoResponse, error)
	GetUserVideosByExercise(context.Context, *GetUserVideosByExerciseRequest) (*GetUserVideosByExerciseResponse, error)
	SaveAnalysis(context.Context, *VideoAnalysisRequest) (*SaveAnalysisResponse, error)
	mustEmbedUnimplementedDBServiceServer()
}

// UnimplementedDBServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDBServiceServer struct{}

func (UnimplementedDBServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedDBServiceServer) SaveUploadedVideo(context.Context, *UploadVideoRequest) (*UploadVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveUploadedVideo not implemented")
}
func (UnimplementedDBServiceServer) GetUserVideosByExercise(context.Context, *GetUserVideosByExerciseRequest) (*GetUserVideosByExerciseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserVideosByExercise not implemented")
}
func (UnimplementedDBServiceServer) SaveAnalysis(context.Context, *VideoAnalysisRequest) (*SaveAnalysisResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveAnalysis not implemented")
}
func (UnimplementedDBServiceServer) mustEmbedUnimplementedDBServiceServer() {}
func (UnimplementedDBServiceServer) testEmbeddedByValue()                   {}

// UnsafeDBServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DBServiceServer will
// result in compilation errors.
type UnsafeDBServiceServer interface {
	mustEmbedUnimplementedDBServiceServer()
}

func RegisterDBServiceServer(s grpc.ServiceRegistrar, srv DBServiceServer) {
	// If the following call pancis, it indicates UnimplementedDBServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DBService_ServiceDesc, srv)
}

func _DBService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_SaveUploadedVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).SaveUploadedVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_SaveUploadedVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).SaveUploadedVideo(ctx, req.(*UploadVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_GetUserVideosByExercise_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserVideosByExerciseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).GetUserVideosByExercise(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_GetUserVideosByExercise_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).GetUserVideosByExercise(ctx, req.(*GetUserVideosByExerciseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_SaveAnalysis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VideoAnalysisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).SaveAnalysis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBService_SaveAnalysis_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).SaveAnalysis(ctx, req.(*VideoAnalysisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DBService_ServiceDesc is the grpc.ServiceDesc for DBService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DBService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "db.DBService",
	HandlerType: (*DBServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _DBService_GetUser_Handler,
		},
		{
			MethodName: "SaveUploadedVideo",
			Handler:    _DBService_SaveUploadedVideo_Handler,
		},
		{
			MethodName: "GetUserVideosByExercise",
			Handler:    _DBService_GetUserVideosByExercise_Handler,
		},
		{
			MethodName: "SaveAnalysis",
			Handler:    _DBService_SaveAnalysis_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/db.proto",
}

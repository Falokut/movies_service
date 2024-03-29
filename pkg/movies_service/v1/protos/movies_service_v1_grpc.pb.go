// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: movies_service_v1.proto

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MoviesServiceV1Client is the client API for MoviesServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MoviesServiceV1Client interface {
	// Returns movie with the specified id.
	GetMovie(ctx context.Context, in *GetMovieRequest, opts ...grpc.CallOption) (*Movie, error)
	// Returns movies previews with the specified filter.
	GetMoviesPreview(ctx context.Context, in *GetMoviesPreviewRequest, opts ...grpc.CallOption) (*MoviesPreview, error)
	// Returns movies previews with the specified ids.
	GetMoviesPreviewByIDs(ctx context.Context, in *GetMoviesPreviewByIDsRequest, opts ...grpc.CallOption) (*MoviesPreview, error)
	// Returns all age ratings.
	GetAgeRatings(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AgeRatings, error)
	// Returns all genres.
	GetGenres(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Genres, error)
	// Returns all countries.
	GetCountries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Countries, error)
}

type moviesServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewMoviesServiceV1Client(cc grpc.ClientConnInterface) MoviesServiceV1Client {
	return &moviesServiceV1Client{cc}
}

func (c *moviesServiceV1Client) GetMovie(ctx context.Context, in *GetMovieRequest, opts ...grpc.CallOption) (*Movie, error) {
	out := new(Movie)
	err := c.cc.Invoke(ctx, "/movies_service.moviesServiceV1/GetMovie", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesServiceV1Client) GetMoviesPreview(ctx context.Context, in *GetMoviesPreviewRequest, opts ...grpc.CallOption) (*MoviesPreview, error) {
	out := new(MoviesPreview)
	err := c.cc.Invoke(ctx, "/movies_service.moviesServiceV1/GetMoviesPreview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesServiceV1Client) GetMoviesPreviewByIDs(ctx context.Context, in *GetMoviesPreviewByIDsRequest, opts ...grpc.CallOption) (*MoviesPreview, error) {
	out := new(MoviesPreview)
	err := c.cc.Invoke(ctx, "/movies_service.moviesServiceV1/GetMoviesPreviewByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesServiceV1Client) GetAgeRatings(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AgeRatings, error) {
	out := new(AgeRatings)
	err := c.cc.Invoke(ctx, "/movies_service.moviesServiceV1/GetAgeRatings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesServiceV1Client) GetGenres(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Genres, error) {
	out := new(Genres)
	err := c.cc.Invoke(ctx, "/movies_service.moviesServiceV1/GetGenres", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesServiceV1Client) GetCountries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Countries, error) {
	out := new(Countries)
	err := c.cc.Invoke(ctx, "/movies_service.moviesServiceV1/GetCountries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MoviesServiceV1Server is the server API for MoviesServiceV1 service.
// All implementations must embed UnimplementedMoviesServiceV1Server
// for forward compatibility
type MoviesServiceV1Server interface {
	// Returns movie with the specified id.
	GetMovie(context.Context, *GetMovieRequest) (*Movie, error)
	// Returns movies previews with the specified filter.
	GetMoviesPreview(context.Context, *GetMoviesPreviewRequest) (*MoviesPreview, error)
	// Returns movies previews with the specified ids.
	GetMoviesPreviewByIDs(context.Context, *GetMoviesPreviewByIDsRequest) (*MoviesPreview, error)
	// Returns all age ratings.
	GetAgeRatings(context.Context, *emptypb.Empty) (*AgeRatings, error)
	// Returns all genres.
	GetGenres(context.Context, *emptypb.Empty) (*Genres, error)
	// Returns all countries.
	GetCountries(context.Context, *emptypb.Empty) (*Countries, error)
	mustEmbedUnimplementedMoviesServiceV1Server()
}

// UnimplementedMoviesServiceV1Server must be embedded to have forward compatible implementations.
type UnimplementedMoviesServiceV1Server struct {
}

func (UnimplementedMoviesServiceV1Server) GetMovie(context.Context, *GetMovieRequest) (*Movie, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovie not implemented")
}
func (UnimplementedMoviesServiceV1Server) GetMoviesPreview(context.Context, *GetMoviesPreviewRequest) (*MoviesPreview, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMoviesPreview not implemented")
}
func (UnimplementedMoviesServiceV1Server) GetMoviesPreviewByIDs(context.Context, *GetMoviesPreviewByIDsRequest) (*MoviesPreview, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMoviesPreviewByIDs not implemented")
}
func (UnimplementedMoviesServiceV1Server) GetAgeRatings(context.Context, *emptypb.Empty) (*AgeRatings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAgeRatings not implemented")
}
func (UnimplementedMoviesServiceV1Server) GetGenres(context.Context, *emptypb.Empty) (*Genres, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGenres not implemented")
}
func (UnimplementedMoviesServiceV1Server) GetCountries(context.Context, *emptypb.Empty) (*Countries, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCountries not implemented")
}
func (UnimplementedMoviesServiceV1Server) mustEmbedUnimplementedMoviesServiceV1Server() {}

// UnsafeMoviesServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MoviesServiceV1Server will
// result in compilation errors.
type UnsafeMoviesServiceV1Server interface {
	mustEmbedUnimplementedMoviesServiceV1Server()
}

func RegisterMoviesServiceV1Server(s grpc.ServiceRegistrar, srv MoviesServiceV1Server) {
	s.RegisterService(&MoviesServiceV1_ServiceDesc, srv)
}

func _MoviesServiceV1_GetMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMovieRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesServiceV1Server).GetMovie(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_service.moviesServiceV1/GetMovie",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesServiceV1Server).GetMovie(ctx, req.(*GetMovieRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesServiceV1_GetMoviesPreview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMoviesPreviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesServiceV1Server).GetMoviesPreview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_service.moviesServiceV1/GetMoviesPreview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesServiceV1Server).GetMoviesPreview(ctx, req.(*GetMoviesPreviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesServiceV1_GetMoviesPreviewByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMoviesPreviewByIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesServiceV1Server).GetMoviesPreviewByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_service.moviesServiceV1/GetMoviesPreviewByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesServiceV1Server).GetMoviesPreviewByIDs(ctx, req.(*GetMoviesPreviewByIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesServiceV1_GetAgeRatings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesServiceV1Server).GetAgeRatings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_service.moviesServiceV1/GetAgeRatings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesServiceV1Server).GetAgeRatings(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesServiceV1_GetGenres_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesServiceV1Server).GetGenres(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_service.moviesServiceV1/GetGenres",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesServiceV1Server).GetGenres(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesServiceV1_GetCountries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesServiceV1Server).GetCountries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_service.moviesServiceV1/GetCountries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesServiceV1Server).GetCountries(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MoviesServiceV1_ServiceDesc is the grpc.ServiceDesc for MoviesServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MoviesServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "movies_service.moviesServiceV1",
	HandlerType: (*MoviesServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMovie",
			Handler:    _MoviesServiceV1_GetMovie_Handler,
		},
		{
			MethodName: "GetMoviesPreview",
			Handler:    _MoviesServiceV1_GetMoviesPreview_Handler,
		},
		{
			MethodName: "GetMoviesPreviewByIDs",
			Handler:    _MoviesServiceV1_GetMoviesPreviewByIDs_Handler,
		},
		{
			MethodName: "GetAgeRatings",
			Handler:    _MoviesServiceV1_GetAgeRatings_Handler,
		},
		{
			MethodName: "GetGenres",
			Handler:    _MoviesServiceV1_GetGenres_Handler,
		},
		{
			MethodName: "GetCountries",
			Handler:    _MoviesServiceV1_GetCountries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "movies_service_v1.proto",
}

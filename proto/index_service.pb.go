// Code generated by protoc-gen-go.
// source: proto/index_service.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	proto/index_service.proto

It has these top-level messages:
	GetIndexInfoRequest
	GetIndexInfoResponse
	PutDocumentRequest
	PutDocumentResponse
	GetDocumentRequest
	GetDocumentResponse
	DeleteDocumentRequest
	DeleteDocumentResponse
	BulkRequest
	BulkResponse
	SearchRequest
	SearchResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/any"

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

type GetIndexInfoRequest struct {
	IndexPath    bool `protobuf:"varint,1,opt,name=index_path,json=indexPath" json:"index_path,omitempty"`
	IndexMapping bool `protobuf:"varint,2,opt,name=index_mapping,json=indexMapping" json:"index_mapping,omitempty"`
	IndexType    bool `protobuf:"varint,3,opt,name=index_type,json=indexType" json:"index_type,omitempty"`
	Kvstore      bool `protobuf:"varint,4,opt,name=kvstore" json:"kvstore,omitempty"`
	Kvconfig     bool `protobuf:"varint,5,opt,name=kvconfig" json:"kvconfig,omitempty"`
}

func (m *GetIndexInfoRequest) Reset()                    { *m = GetIndexInfoRequest{} }
func (m *GetIndexInfoRequest) String() string            { return proto1.CompactTextString(m) }
func (*GetIndexInfoRequest) ProtoMessage()               {}
func (*GetIndexInfoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetIndexInfoRequest) GetIndexPath() bool {
	if m != nil {
		return m.IndexPath
	}
	return false
}

func (m *GetIndexInfoRequest) GetIndexMapping() bool {
	if m != nil {
		return m.IndexMapping
	}
	return false
}

func (m *GetIndexInfoRequest) GetIndexType() bool {
	if m != nil {
		return m.IndexType
	}
	return false
}

func (m *GetIndexInfoRequest) GetKvstore() bool {
	if m != nil {
		return m.Kvstore
	}
	return false
}

func (m *GetIndexInfoRequest) GetKvconfig() bool {
	if m != nil {
		return m.Kvconfig
	}
	return false
}

type GetIndexInfoResponse struct {
	IndexPath    string               `protobuf:"bytes,1,opt,name=index_path,json=indexPath" json:"index_path,omitempty"`
	IndexMapping *google_protobuf.Any `protobuf:"bytes,2,opt,name=index_mapping,json=indexMapping" json:"index_mapping,omitempty"`
	IndexType    string               `protobuf:"bytes,3,opt,name=index_type,json=indexType" json:"index_type,omitempty"`
	Kvstore      string               `protobuf:"bytes,4,opt,name=kvstore" json:"kvstore,omitempty"`
	Kvconfig     *google_protobuf.Any `protobuf:"bytes,5,opt,name=kvconfig" json:"kvconfig,omitempty"`
}

func (m *GetIndexInfoResponse) Reset()                    { *m = GetIndexInfoResponse{} }
func (m *GetIndexInfoResponse) String() string            { return proto1.CompactTextString(m) }
func (*GetIndexInfoResponse) ProtoMessage()               {}
func (*GetIndexInfoResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetIndexInfoResponse) GetIndexPath() string {
	if m != nil {
		return m.IndexPath
	}
	return ""
}

func (m *GetIndexInfoResponse) GetIndexMapping() *google_protobuf.Any {
	if m != nil {
		return m.IndexMapping
	}
	return nil
}

func (m *GetIndexInfoResponse) GetIndexType() string {
	if m != nil {
		return m.IndexType
	}
	return ""
}

func (m *GetIndexInfoResponse) GetKvstore() string {
	if m != nil {
		return m.Kvstore
	}
	return ""
}

func (m *GetIndexInfoResponse) GetKvconfig() *google_protobuf.Any {
	if m != nil {
		return m.Kvconfig
	}
	return nil
}

type PutDocumentRequest struct {
	Id     string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Fields *google_protobuf.Any `protobuf:"bytes,2,opt,name=fields" json:"fields,omitempty"`
}

func (m *PutDocumentRequest) Reset()                    { *m = PutDocumentRequest{} }
func (m *PutDocumentRequest) String() string            { return proto1.CompactTextString(m) }
func (*PutDocumentRequest) ProtoMessage()               {}
func (*PutDocumentRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PutDocumentRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PutDocumentRequest) GetFields() *google_protobuf.Any {
	if m != nil {
		return m.Fields
	}
	return nil
}

type PutDocumentResponse struct {
	Id     string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Fields *google_protobuf.Any `protobuf:"bytes,2,opt,name=fields" json:"fields,omitempty"`
}

func (m *PutDocumentResponse) Reset()                    { *m = PutDocumentResponse{} }
func (m *PutDocumentResponse) String() string            { return proto1.CompactTextString(m) }
func (*PutDocumentResponse) ProtoMessage()               {}
func (*PutDocumentResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PutDocumentResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PutDocumentResponse) GetFields() *google_protobuf.Any {
	if m != nil {
		return m.Fields
	}
	return nil
}

type GetDocumentRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *GetDocumentRequest) Reset()                    { *m = GetDocumentRequest{} }
func (m *GetDocumentRequest) String() string            { return proto1.CompactTextString(m) }
func (*GetDocumentRequest) ProtoMessage()               {}
func (*GetDocumentRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetDocumentRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetDocumentResponse struct {
	Id     string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Fields *google_protobuf.Any `protobuf:"bytes,2,opt,name=fields" json:"fields,omitempty"`
}

func (m *GetDocumentResponse) Reset()                    { *m = GetDocumentResponse{} }
func (m *GetDocumentResponse) String() string            { return proto1.CompactTextString(m) }
func (*GetDocumentResponse) ProtoMessage()               {}
func (*GetDocumentResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GetDocumentResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GetDocumentResponse) GetFields() *google_protobuf.Any {
	if m != nil {
		return m.Fields
	}
	return nil
}

type DeleteDocumentRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteDocumentRequest) Reset()                    { *m = DeleteDocumentRequest{} }
func (m *DeleteDocumentRequest) String() string            { return proto1.CompactTextString(m) }
func (*DeleteDocumentRequest) ProtoMessage()               {}
func (*DeleteDocumentRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *DeleteDocumentRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteDocumentResponse struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteDocumentResponse) Reset()                    { *m = DeleteDocumentResponse{} }
func (m *DeleteDocumentResponse) String() string            { return proto1.CompactTextString(m) }
func (*DeleteDocumentResponse) ProtoMessage()               {}
func (*DeleteDocumentResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *DeleteDocumentResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type BulkRequest struct {
	BatchSize      int32                        `protobuf:"varint,1,opt,name=batch_size,json=batchSize" json:"batch_size,omitempty"`
	UpdateRequests []*BulkRequest_UpdateRequest `protobuf:"bytes,2,rep,name=update_requests,json=updateRequests" json:"update_requests,omitempty"`
}

func (m *BulkRequest) Reset()                    { *m = BulkRequest{} }
func (m *BulkRequest) String() string            { return proto1.CompactTextString(m) }
func (*BulkRequest) ProtoMessage()               {}
func (*BulkRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *BulkRequest) GetBatchSize() int32 {
	if m != nil {
		return m.BatchSize
	}
	return 0
}

func (m *BulkRequest) GetUpdateRequests() []*BulkRequest_UpdateRequest {
	if m != nil {
		return m.UpdateRequests
	}
	return nil
}

type BulkRequest_UpdateRequest struct {
	Method   string                              `protobuf:"bytes,1,opt,name=method" json:"method,omitempty"`
	Document *BulkRequest_UpdateRequest_Document `protobuf:"bytes,2,opt,name=document" json:"document,omitempty"`
}

func (m *BulkRequest_UpdateRequest) Reset()                    { *m = BulkRequest_UpdateRequest{} }
func (m *BulkRequest_UpdateRequest) String() string            { return proto1.CompactTextString(m) }
func (*BulkRequest_UpdateRequest) ProtoMessage()               {}
func (*BulkRequest_UpdateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8, 0} }

func (m *BulkRequest_UpdateRequest) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *BulkRequest_UpdateRequest) GetDocument() *BulkRequest_UpdateRequest_Document {
	if m != nil {
		return m.Document
	}
	return nil
}

type BulkRequest_UpdateRequest_Document struct {
	Id     string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Fields *google_protobuf.Any `protobuf:"bytes,2,opt,name=fields" json:"fields,omitempty"`
}

func (m *BulkRequest_UpdateRequest_Document) Reset()         { *m = BulkRequest_UpdateRequest_Document{} }
func (m *BulkRequest_UpdateRequest_Document) String() string { return proto1.CompactTextString(m) }
func (*BulkRequest_UpdateRequest_Document) ProtoMessage()    {}
func (*BulkRequest_UpdateRequest_Document) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{8, 0, 0}
}

func (m *BulkRequest_UpdateRequest_Document) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BulkRequest_UpdateRequest_Document) GetFields() *google_protobuf.Any {
	if m != nil {
		return m.Fields
	}
	return nil
}

type BulkResponse struct {
	PutCount         int32 `protobuf:"varint,1,opt,name=put_count,json=putCount" json:"put_count,omitempty"`
	PutErrorCount    int32 `protobuf:"varint,2,opt,name=put_error_count,json=putErrorCount" json:"put_error_count,omitempty"`
	DeleteCount      int32 `protobuf:"varint,3,opt,name=delete_count,json=deleteCount" json:"delete_count,omitempty"`
	MethodErrorCount int32 `protobuf:"varint,4,opt,name=method_error_count,json=methodErrorCount" json:"method_error_count,omitempty"`
}

func (m *BulkResponse) Reset()                    { *m = BulkResponse{} }
func (m *BulkResponse) String() string            { return proto1.CompactTextString(m) }
func (*BulkResponse) ProtoMessage()               {}
func (*BulkResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *BulkResponse) GetPutCount() int32 {
	if m != nil {
		return m.PutCount
	}
	return 0
}

func (m *BulkResponse) GetPutErrorCount() int32 {
	if m != nil {
		return m.PutErrorCount
	}
	return 0
}

func (m *BulkResponse) GetDeleteCount() int32 {
	if m != nil {
		return m.DeleteCount
	}
	return 0
}

func (m *BulkResponse) GetMethodErrorCount() int32 {
	if m != nil {
		return m.MethodErrorCount
	}
	return 0
}

type SearchRequest struct {
	SearchRequest *google_protobuf.Any `protobuf:"bytes,1,opt,name=search_request,json=searchRequest" json:"search_request,omitempty"`
}

func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
func (m *SearchRequest) String() string            { return proto1.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *SearchRequest) GetSearchRequest() *google_protobuf.Any {
	if m != nil {
		return m.SearchRequest
	}
	return nil
}

type SearchResponse struct {
	SearchResult *google_protobuf.Any `protobuf:"bytes,1,opt,name=search_result,json=searchResult" json:"search_result,omitempty"`
}

func (m *SearchResponse) Reset()                    { *m = SearchResponse{} }
func (m *SearchResponse) String() string            { return proto1.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *SearchResponse) GetSearchResult() *google_protobuf.Any {
	if m != nil {
		return m.SearchResult
	}
	return nil
}

func init() {
	proto1.RegisterType((*GetIndexInfoRequest)(nil), "proto.GetIndexInfoRequest")
	proto1.RegisterType((*GetIndexInfoResponse)(nil), "proto.GetIndexInfoResponse")
	proto1.RegisterType((*PutDocumentRequest)(nil), "proto.PutDocumentRequest")
	proto1.RegisterType((*PutDocumentResponse)(nil), "proto.PutDocumentResponse")
	proto1.RegisterType((*GetDocumentRequest)(nil), "proto.GetDocumentRequest")
	proto1.RegisterType((*GetDocumentResponse)(nil), "proto.GetDocumentResponse")
	proto1.RegisterType((*DeleteDocumentRequest)(nil), "proto.DeleteDocumentRequest")
	proto1.RegisterType((*DeleteDocumentResponse)(nil), "proto.DeleteDocumentResponse")
	proto1.RegisterType((*BulkRequest)(nil), "proto.BulkRequest")
	proto1.RegisterType((*BulkRequest_UpdateRequest)(nil), "proto.BulkRequest.UpdateRequest")
	proto1.RegisterType((*BulkRequest_UpdateRequest_Document)(nil), "proto.BulkRequest.UpdateRequest.Document")
	proto1.RegisterType((*BulkResponse)(nil), "proto.BulkResponse")
	proto1.RegisterType((*SearchRequest)(nil), "proto.SearchRequest")
	proto1.RegisterType((*SearchResponse)(nil), "proto.SearchResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Index service

type IndexClient interface {
	GetIndexInfo(ctx context.Context, in *GetIndexInfoRequest, opts ...grpc.CallOption) (*GetIndexInfoResponse, error)
	PutDocument(ctx context.Context, in *PutDocumentRequest, opts ...grpc.CallOption) (*PutDocumentResponse, error)
	GetDocument(ctx context.Context, in *GetDocumentRequest, opts ...grpc.CallOption) (*GetDocumentResponse, error)
	DeleteDocument(ctx context.Context, in *DeleteDocumentRequest, opts ...grpc.CallOption) (*DeleteDocumentResponse, error)
	Bulk(ctx context.Context, in *BulkRequest, opts ...grpc.CallOption) (*BulkResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
}

type indexClient struct {
	cc *grpc.ClientConn
}

func NewIndexClient(cc *grpc.ClientConn) IndexClient {
	return &indexClient{cc}
}

func (c *indexClient) GetIndexInfo(ctx context.Context, in *GetIndexInfoRequest, opts ...grpc.CallOption) (*GetIndexInfoResponse, error) {
	out := new(GetIndexInfoResponse)
	err := grpc.Invoke(ctx, "/proto.Index/GetIndexInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) PutDocument(ctx context.Context, in *PutDocumentRequest, opts ...grpc.CallOption) (*PutDocumentResponse, error) {
	out := new(PutDocumentResponse)
	err := grpc.Invoke(ctx, "/proto.Index/PutDocument", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) GetDocument(ctx context.Context, in *GetDocumentRequest, opts ...grpc.CallOption) (*GetDocumentResponse, error) {
	out := new(GetDocumentResponse)
	err := grpc.Invoke(ctx, "/proto.Index/GetDocument", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) DeleteDocument(ctx context.Context, in *DeleteDocumentRequest, opts ...grpc.CallOption) (*DeleteDocumentResponse, error) {
	out := new(DeleteDocumentResponse)
	err := grpc.Invoke(ctx, "/proto.Index/DeleteDocument", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) Bulk(ctx context.Context, in *BulkRequest, opts ...grpc.CallOption) (*BulkResponse, error) {
	out := new(BulkResponse)
	err := grpc.Invoke(ctx, "/proto.Index/Bulk", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := grpc.Invoke(ctx, "/proto.Index/Search", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Index service

type IndexServer interface {
	GetIndexInfo(context.Context, *GetIndexInfoRequest) (*GetIndexInfoResponse, error)
	PutDocument(context.Context, *PutDocumentRequest) (*PutDocumentResponse, error)
	GetDocument(context.Context, *GetDocumentRequest) (*GetDocumentResponse, error)
	DeleteDocument(context.Context, *DeleteDocumentRequest) (*DeleteDocumentResponse, error)
	Bulk(context.Context, *BulkRequest) (*BulkResponse, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
}

func RegisterIndexServer(s *grpc.Server, srv IndexServer) {
	s.RegisterService(&_Index_serviceDesc, srv)
}

func _Index_GetIndexInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIndexInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).GetIndexInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Index/GetIndexInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).GetIndexInfo(ctx, req.(*GetIndexInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_PutDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutDocumentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).PutDocument(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Index/PutDocument",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).PutDocument(ctx, req.(*PutDocumentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_GetDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDocumentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).GetDocument(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Index/GetDocument",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).GetDocument(ctx, req.(*GetDocumentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_DeleteDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDocumentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).DeleteDocument(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Index/DeleteDocument",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).DeleteDocument(ctx, req.(*DeleteDocumentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_Bulk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).Bulk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Index/Bulk",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).Bulk(ctx, req.(*BulkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Index/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Index_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Index",
	HandlerType: (*IndexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetIndexInfo",
			Handler:    _Index_GetIndexInfo_Handler,
		},
		{
			MethodName: "PutDocument",
			Handler:    _Index_PutDocument_Handler,
		},
		{
			MethodName: "GetDocument",
			Handler:    _Index_GetDocument_Handler,
		},
		{
			MethodName: "DeleteDocument",
			Handler:    _Index_DeleteDocument_Handler,
		},
		{
			MethodName: "Bulk",
			Handler:    _Index_Bulk_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Index_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/index_service.proto",
}

func init() { proto1.RegisterFile("proto/index_service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 656 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0x49, 0x13, 0xe2, 0x49, 0x9c, 0xa2, 0xed, 0x87, 0xdc, 0x2d, 0x95, 0x8a, 0x41, 0x50,
	0xa4, 0xca, 0x85, 0x72, 0x40, 0x15, 0x27, 0xa0, 0xa5, 0x44, 0x80, 0xa8, 0x52, 0x38, 0x47, 0xae,
	0xbd, 0x49, 0xac, 0xa6, 0x5e, 0x63, 0xef, 0x56, 0xa4, 0x7f, 0x86, 0x0b, 0x67, 0x0e, 0x9c, 0xf8,
	0x2f, 0xfc, 0x19, 0xe4, 0xdd, 0xb5, 0x63, 0x27, 0x4e, 0x72, 0x28, 0x27, 0x6b, 0xde, 0xbc, 0x7d,
	0x33, 0x6f, 0xc6, 0xbb, 0xb0, 0x15, 0x46, 0x94, 0xd1, 0x03, 0x3f, 0xf0, 0xc8, 0xf7, 0x5e, 0x4c,
	0xa2, 0x6b, 0xdf, 0x25, 0xb6, 0xc0, 0x50, 0x4d, 0x7c, 0xf0, 0xd6, 0x80, 0xd2, 0xc1, 0x88, 0x1c,
	0x88, 0xe8, 0x82, 0xf7, 0x0f, 0x9c, 0x60, 0x2c, 0x19, 0xd6, 0x2f, 0x0d, 0xd6, 0x4e, 0x09, 0xeb,
	0x24, 0x87, 0x3b, 0x41, 0x9f, 0x76, 0xc9, 0x37, 0x4e, 0x62, 0x86, 0x76, 0x00, 0xa4, 0x60, 0xe8,
	0xb0, 0xa1, 0xa9, 0xed, 0x6a, 0x7b, 0x8d, 0xae, 0x2e, 0x90, 0x33, 0x87, 0x0d, 0xd1, 0x43, 0x30,
	0x64, 0xfa, 0xca, 0x09, 0x43, 0x3f, 0x18, 0x98, 0x15, 0xc1, 0x68, 0x09, 0xf0, 0x93, 0xc4, 0x26,
	0x1a, 0x6c, 0x1c, 0x12, 0xb3, 0x9a, 0xd3, 0xf8, 0x32, 0x0e, 0x09, 0x32, 0xe1, 0xee, 0xe5, 0x75,
	0xcc, 0x68, 0x44, 0xcc, 0x15, 0x91, 0x4b, 0x43, 0x84, 0xa1, 0x71, 0x79, 0xed, 0xd2, 0xa0, 0xef,
	0x0f, 0xcc, 0x9a, 0x48, 0x65, 0xb1, 0xf5, 0x57, 0x83, 0xf5, 0x62, 0xc3, 0x71, 0x48, 0x83, 0x98,
	0x94, 0x74, 0xac, 0xe7, 0x3b, 0x3e, 0x2a, 0xeb, 0xb8, 0x79, 0xb8, 0x6e, 0xcb, 0xd9, 0xd8, 0xe9,
	0x6c, 0xec, 0xd7, 0xc1, 0x78, 0xa9, 0x0f, 0x7d, 0x81, 0x0f, 0x7d, 0xe2, 0xe3, 0xd9, 0x94, 0x8f,
	0x79, 0xe5, 0x26, 0xee, 0xba, 0x80, 0xce, 0x38, 0x3b, 0xa6, 0x2e, 0xbf, 0x22, 0x01, 0x4b, 0x97,
	0xd1, 0x86, 0x8a, 0xef, 0x29, 0x4b, 0x15, 0xdf, 0x43, 0xfb, 0x50, 0xef, 0xfb, 0x64, 0xe4, 0xc5,
	0x0b, 0x4d, 0x28, 0x8e, 0x75, 0x0e, 0x6b, 0x05, 0x4d, 0x35, 0xaf, 0xdb, 0x89, 0x3e, 0x02, 0x74,
	0x4a, 0x96, 0x35, 0x9a, 0x94, 0x2e, 0xb0, 0xfe, 0x4b, 0xe9, 0x27, 0xb0, 0x71, 0x4c, 0x46, 0x84,
	0x91, 0x65, 0xd5, 0xf7, 0x60, 0x73, 0x9a, 0x58, 0xde, 0x80, 0xf5, 0xbb, 0x02, 0xcd, 0x37, 0x7c,
	0x74, 0x99, 0xfb, 0xfb, 0x2f, 0x1c, 0xe6, 0x0e, 0x7b, 0xb1, 0x7f, 0x43, 0x04, 0xaf, 0xd6, 0xd5,
	0x05, 0x72, 0xee, 0xdf, 0x10, 0xd4, 0x81, 0x55, 0x1e, 0x7a, 0x0e, 0x23, 0xbd, 0x48, 0x1e, 0x48,
	0x1a, 0xaf, 0xee, 0x35, 0x0f, 0x77, 0x65, 0xc7, 0x76, 0x4e, 0xcb, 0xfe, 0x2a, 0x98, 0x2a, 0xea,
	0xb6, 0x79, 0x3e, 0x8c, 0xf1, 0x1f, 0x0d, 0x8c, 0x02, 0x03, 0x6d, 0x42, 0xfd, 0x8a, 0xb0, 0x21,
	0x4d, 0xfb, 0x53, 0x11, 0x3a, 0x81, 0x86, 0xa7, 0x7c, 0xa8, 0x31, 0x3d, 0x5d, 0x56, 0xcd, 0xce,
	0x8c, 0x67, 0x47, 0xf1, 0x7b, 0x68, 0xa4, 0xe8, 0x2d, 0xf7, 0xf0, 0x53, 0x83, 0x96, 0x2c, 0xad,
	0xa6, 0xba, 0x0d, 0x7a, 0xc8, 0x59, 0xcf, 0xa5, 0x3c, 0x60, 0x6a, 0x68, 0x8d, 0x90, 0xb3, 0xb7,
	0x49, 0x8c, 0x1e, 0xc3, 0x6a, 0x92, 0x24, 0x51, 0x44, 0x23, 0x45, 0xa9, 0x08, 0x8a, 0x11, 0x72,
	0x76, 0x92, 0xa0, 0x92, 0xf7, 0x00, 0x5a, 0x9e, 0x58, 0x9a, 0x22, 0x55, 0x05, 0xa9, 0x29, 0x31,
	0x49, 0xd9, 0x07, 0x24, 0x67, 0x52, 0x50, 0x5b, 0x11, 0xc4, 0x7b, 0x32, 0x33, 0x11, 0xb4, 0x3e,
	0x82, 0x71, 0x4e, 0x9c, 0xc8, 0x1d, 0xa6, 0x03, 0x7e, 0x05, 0xed, 0x58, 0x00, 0xe9, 0xf6, 0x44,
	0xaf, 0xf3, 0xdc, 0x1a, 0x71, 0xfe, 0xb0, 0xf5, 0x01, 0xda, 0xa9, 0x9a, 0x72, 0x7d, 0x04, 0x46,
	0x26, 0x17, 0xf3, 0xd1, 0x62, 0xb5, 0x56, 0xaa, 0x96, 0x30, 0x0f, 0x7f, 0x54, 0xa1, 0x26, 0x1e,
	0x32, 0xd4, 0x81, 0x56, 0xfe, 0x51, 0x43, 0x58, 0xad, 0xb6, 0xe4, 0x69, 0xc6, 0xdb, 0xa5, 0x39,
	0xd9, 0x8d, 0x75, 0x07, 0xbd, 0x83, 0x66, 0xee, 0xba, 0xa3, 0x2d, 0xc5, 0x9e, 0x7d, 0x56, 0x30,
	0x2e, 0x4b, 0xe5, 0x75, 0x72, 0x77, 0x37, 0xd3, 0x99, 0xbd, 0xf5, 0x18, 0x97, 0xa5, 0x32, 0x9d,
	0xcf, 0xd0, 0x2e, 0xde, 0x42, 0x74, 0x5f, 0xf1, 0x4b, 0x6f, 0x31, 0xde, 0x99, 0x93, 0xcd, 0x04,
	0x9f, 0xc3, 0x4a, 0xf2, 0xdb, 0x21, 0x34, 0xfb, 0xfb, 0xe3, 0xb5, 0x02, 0x96, 0x1d, 0x79, 0x09,
	0x75, 0xb9, 0x35, 0xb4, 0xae, 0x08, 0x85, 0x5f, 0x02, 0x6f, 0x4c, 0xa1, 0xe9, 0xc1, 0x8b, 0xba,
	0xc0, 0x5f, 0xfc, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x87, 0xa2, 0xa1, 0x07, 0x64, 0x07, 0x00, 0x00,
}

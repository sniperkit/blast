// Code generated by protoc-gen-go.
// source: protobuf/message.proto
// DO NOT EDIT!

package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/any"
import _ "github.com/golang/protobuf/ptypes/empty"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PutDocumentRequest struct {
	Id     string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Fields *google_protobuf.Any `protobuf:"bytes,2,opt,name=fields" json:"fields,omitempty"`
}

func (m *PutDocumentRequest) Reset()                    { *m = PutDocumentRequest{} }
func (m *PutDocumentRequest) String() string            { return proto.CompactTextString(m) }
func (*PutDocumentRequest) ProtoMessage()               {}
func (*PutDocumentRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

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
func (m *PutDocumentResponse) String() string            { return proto.CompactTextString(m) }
func (*PutDocumentResponse) ProtoMessage()               {}
func (*PutDocumentResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

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
func (m *GetDocumentRequest) String() string            { return proto.CompactTextString(m) }
func (*GetDocumentRequest) ProtoMessage()               {}
func (*GetDocumentRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

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
func (m *GetDocumentResponse) String() string            { return proto.CompactTextString(m) }
func (*GetDocumentResponse) ProtoMessage()               {}
func (*GetDocumentResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

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
func (m *DeleteDocumentRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteDocumentRequest) ProtoMessage()               {}
func (*DeleteDocumentRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

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
func (m *DeleteDocumentResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteDocumentResponse) ProtoMessage()               {}
func (*DeleteDocumentResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

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
func (m *BulkRequest) String() string            { return proto.CompactTextString(m) }
func (*BulkRequest) ProtoMessage()               {}
func (*BulkRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6} }

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
func (m *BulkRequest_UpdateRequest) String() string            { return proto.CompactTextString(m) }
func (*BulkRequest_UpdateRequest) ProtoMessage()               {}
func (*BulkRequest_UpdateRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{6, 0} }

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
func (m *BulkRequest_UpdateRequest_Document) String() string { return proto.CompactTextString(m) }
func (*BulkRequest_UpdateRequest_Document) ProtoMessage()    {}
func (*BulkRequest_UpdateRequest_Document) Descriptor() ([]byte, []int) {
	return fileDescriptor2, []int{6, 0, 0}
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
func (m *BulkResponse) String() string            { return proto.CompactTextString(m) }
func (*BulkResponse) ProtoMessage()               {}
func (*BulkResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{7} }

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
func (m *SearchRequest) String() string            { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()               {}
func (*SearchRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{8} }

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
func (m *SearchResponse) String() string            { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()               {}
func (*SearchResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{9} }

func (m *SearchResponse) GetSearchResult() *google_protobuf.Any {
	if m != nil {
		return m.SearchResult
	}
	return nil
}

type NodeInfo struct {
	Cluster   string `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
	Node      string `protobuf:"bytes,2,opt,name=node" json:"node,omitempty"`
	Role      string `protobuf:"bytes,3,opt,name=role" json:"role,omitempty"`
	Status    string `protobuf:"bytes,4,opt,name=status" json:"status,omitempty"`
	Timestamp string `protobuf:"bytes,5,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *NodeInfo) Reset()                    { *m = NodeInfo{} }
func (m *NodeInfo) String() string            { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()               {}
func (*NodeInfo) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{10} }

func (m *NodeInfo) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *NodeInfo) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *NodeInfo) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

func (m *NodeInfo) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *NodeInfo) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type PutNodeRequest struct {
	Cluster  string    `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
	Node     string    `protobuf:"bytes,2,opt,name=node" json:"node,omitempty"`
	NodeInfo *NodeInfo `protobuf:"bytes,3,opt,name=node_info,json=nodeInfo" json:"node_info,omitempty"`
}

func (m *PutNodeRequest) Reset()                    { *m = PutNodeRequest{} }
func (m *PutNodeRequest) String() string            { return proto.CompactTextString(m) }
func (*PutNodeRequest) ProtoMessage()               {}
func (*PutNodeRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{11} }

func (m *PutNodeRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *PutNodeRequest) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *PutNodeRequest) GetNodeInfo() *NodeInfo {
	if m != nil {
		return m.NodeInfo
	}
	return nil
}

type GetNodeRequest struct {
	Cluster string `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
	Node    string `protobuf:"bytes,2,opt,name=node" json:"node,omitempty"`
}

func (m *GetNodeRequest) Reset()                    { *m = GetNodeRequest{} }
func (m *GetNodeRequest) String() string            { return proto.CompactTextString(m) }
func (*GetNodeRequest) ProtoMessage()               {}
func (*GetNodeRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{12} }

func (m *GetNodeRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *GetNodeRequest) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

type GetNodeResponse struct {
	NodeInfo *NodeInfo `protobuf:"bytes,1,opt,name=node_info,json=nodeInfo" json:"node_info,omitempty"`
}

func (m *GetNodeResponse) Reset()                    { *m = GetNodeResponse{} }
func (m *GetNodeResponse) String() string            { return proto.CompactTextString(m) }
func (*GetNodeResponse) ProtoMessage()               {}
func (*GetNodeResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{13} }

func (m *GetNodeResponse) GetNodeInfo() *NodeInfo {
	if m != nil {
		return m.NodeInfo
	}
	return nil
}

type DeleteNodeRequest struct {
	Cluster string `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
	Node    string `protobuf:"bytes,2,opt,name=node" json:"node,omitempty"`
}

func (m *DeleteNodeRequest) Reset()                    { *m = DeleteNodeRequest{} }
func (m *DeleteNodeRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteNodeRequest) ProtoMessage()               {}
func (*DeleteNodeRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{14} }

func (m *DeleteNodeRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *DeleteNodeRequest) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

type PutIndexMappingRequest struct {
	Cluster      string               `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
	IndexMapping *google_protobuf.Any `protobuf:"bytes,2,opt,name=index_mapping,json=indexMapping" json:"index_mapping,omitempty"`
}

func (m *PutIndexMappingRequest) Reset()                    { *m = PutIndexMappingRequest{} }
func (m *PutIndexMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*PutIndexMappingRequest) ProtoMessage()               {}
func (*PutIndexMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{15} }

func (m *PutIndexMappingRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *PutIndexMappingRequest) GetIndexMapping() *google_protobuf.Any {
	if m != nil {
		return m.IndexMapping
	}
	return nil
}

type GetIndexMappingRequest struct {
	Cluster string `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
}

func (m *GetIndexMappingRequest) Reset()                    { *m = GetIndexMappingRequest{} }
func (m *GetIndexMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*GetIndexMappingRequest) ProtoMessage()               {}
func (*GetIndexMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{16} }

func (m *GetIndexMappingRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

type GetIndexMappingResponse struct {
	IndexMapping *google_protobuf.Any `protobuf:"bytes,1,opt,name=index_mapping,json=indexMapping" json:"index_mapping,omitempty"`
}

func (m *GetIndexMappingResponse) Reset()                    { *m = GetIndexMappingResponse{} }
func (m *GetIndexMappingResponse) String() string            { return proto.CompactTextString(m) }
func (*GetIndexMappingResponse) ProtoMessage()               {}
func (*GetIndexMappingResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{17} }

func (m *GetIndexMappingResponse) GetIndexMapping() *google_protobuf.Any {
	if m != nil {
		return m.IndexMapping
	}
	return nil
}

type DeleteIndexMappingRequest struct {
	Cluster string `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
}

func (m *DeleteIndexMappingRequest) Reset()                    { *m = DeleteIndexMappingRequest{} }
func (m *DeleteIndexMappingRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexMappingRequest) ProtoMessage()               {}
func (*DeleteIndexMappingRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{18} }

func (m *DeleteIndexMappingRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

type PutIndexMetaRequest struct {
	Cluster   string               `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
	IndexMeta *google_protobuf.Any `protobuf:"bytes,2,opt,name=index_meta,json=indexMeta" json:"index_meta,omitempty"`
}

func (m *PutIndexMetaRequest) Reset()                    { *m = PutIndexMetaRequest{} }
func (m *PutIndexMetaRequest) String() string            { return proto.CompactTextString(m) }
func (*PutIndexMetaRequest) ProtoMessage()               {}
func (*PutIndexMetaRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{19} }

func (m *PutIndexMetaRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *PutIndexMetaRequest) GetIndexMeta() *google_protobuf.Any {
	if m != nil {
		return m.IndexMeta
	}
	return nil
}

type GetIndexMetaRequest struct {
	Cluster string `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
}

func (m *GetIndexMetaRequest) Reset()                    { *m = GetIndexMetaRequest{} }
func (m *GetIndexMetaRequest) String() string            { return proto.CompactTextString(m) }
func (*GetIndexMetaRequest) ProtoMessage()               {}
func (*GetIndexMetaRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{20} }

func (m *GetIndexMetaRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

type GetIndexMetaResponse struct {
	IndexMeta *google_protobuf.Any `protobuf:"bytes,1,opt,name=index_meta,json=indexMeta" json:"index_meta,omitempty"`
}

func (m *GetIndexMetaResponse) Reset()                    { *m = GetIndexMetaResponse{} }
func (m *GetIndexMetaResponse) String() string            { return proto.CompactTextString(m) }
func (*GetIndexMetaResponse) ProtoMessage()               {}
func (*GetIndexMetaResponse) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{21} }

func (m *GetIndexMetaResponse) GetIndexMeta() *google_protobuf.Any {
	if m != nil {
		return m.IndexMeta
	}
	return nil
}

type DeleteIndexMetaRequest struct {
	Cluster string `protobuf:"bytes,1,opt,name=cluster" json:"cluster,omitempty"`
}

func (m *DeleteIndexMetaRequest) Reset()                    { *m = DeleteIndexMetaRequest{} }
func (m *DeleteIndexMetaRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteIndexMetaRequest) ProtoMessage()               {}
func (*DeleteIndexMetaRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{22} }

func (m *DeleteIndexMetaRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func init() {
	proto.RegisterType((*PutDocumentRequest)(nil), "blast.PutDocumentRequest")
	proto.RegisterType((*PutDocumentResponse)(nil), "blast.PutDocumentResponse")
	proto.RegisterType((*GetDocumentRequest)(nil), "blast.GetDocumentRequest")
	proto.RegisterType((*GetDocumentResponse)(nil), "blast.GetDocumentResponse")
	proto.RegisterType((*DeleteDocumentRequest)(nil), "blast.DeleteDocumentRequest")
	proto.RegisterType((*DeleteDocumentResponse)(nil), "blast.DeleteDocumentResponse")
	proto.RegisterType((*BulkRequest)(nil), "blast.BulkRequest")
	proto.RegisterType((*BulkRequest_UpdateRequest)(nil), "blast.BulkRequest.UpdateRequest")
	proto.RegisterType((*BulkRequest_UpdateRequest_Document)(nil), "blast.BulkRequest.UpdateRequest.Document")
	proto.RegisterType((*BulkResponse)(nil), "blast.BulkResponse")
	proto.RegisterType((*SearchRequest)(nil), "blast.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "blast.SearchResponse")
	proto.RegisterType((*NodeInfo)(nil), "blast.NodeInfo")
	proto.RegisterType((*PutNodeRequest)(nil), "blast.PutNodeRequest")
	proto.RegisterType((*GetNodeRequest)(nil), "blast.GetNodeRequest")
	proto.RegisterType((*GetNodeResponse)(nil), "blast.GetNodeResponse")
	proto.RegisterType((*DeleteNodeRequest)(nil), "blast.DeleteNodeRequest")
	proto.RegisterType((*PutIndexMappingRequest)(nil), "blast.PutIndexMappingRequest")
	proto.RegisterType((*GetIndexMappingRequest)(nil), "blast.GetIndexMappingRequest")
	proto.RegisterType((*GetIndexMappingResponse)(nil), "blast.GetIndexMappingResponse")
	proto.RegisterType((*DeleteIndexMappingRequest)(nil), "blast.DeleteIndexMappingRequest")
	proto.RegisterType((*PutIndexMetaRequest)(nil), "blast.PutIndexMetaRequest")
	proto.RegisterType((*GetIndexMetaRequest)(nil), "blast.GetIndexMetaRequest")
	proto.RegisterType((*GetIndexMetaResponse)(nil), "blast.GetIndexMetaResponse")
	proto.RegisterType((*DeleteIndexMetaRequest)(nil), "blast.DeleteIndexMetaRequest")
}

func init() { proto.RegisterFile("protobuf/message.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 681 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x55, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0x95, 0xd3, 0x9f, 0x2f, 0x9e, 0xfc, 0x7d, 0x2c, 0x25, 0xa4, 0x29, 0x48, 0xc5, 0x42, 0x10,
	0xa4, 0xca, 0x95, 0x52, 0x71, 0x51, 0x21, 0x81, 0x5a, 0x5a, 0x95, 0xaa, 0x80, 0x2a, 0x17, 0x6e,
	0xb8, 0x89, 0x9c, 0x78, 0x92, 0x5a, 0xd8, 0x5e, 0xe3, 0xdd, 0x95, 0x68, 0x2f, 0x79, 0x16, 0x9e,
	0x80, 0x2b, 0x1e, 0x0f, 0x79, 0x7f, 0xd2, 0xb8, 0x40, 0x9d, 0xb6, 0x5c, 0x65, 0x67, 0xf6, 0xcc,
	0x99, 0xb3, 0x67, 0x36, 0x5e, 0x68, 0xa7, 0x19, 0xe5, 0x74, 0x28, 0xc6, 0x9b, 0x31, 0x32, 0xe6,
	0x4f, 0xd0, 0x95, 0x09, 0xb2, 0x34, 0x8c, 0x7c, 0xc6, 0xbb, 0xab, 0x13, 0x4a, 0x27, 0x11, 0x6e,
	0x4e, 0x51, 0x7e, 0x72, 0xa6, 0x10, 0xdd, 0xb5, 0xcb, 0x5b, 0x18, 0xa7, 0x5c, 0x6f, 0x3a, 0x1e,
	0x90, 0x63, 0xc1, 0xf7, 0xe8, 0x48, 0xc4, 0x98, 0x70, 0x0f, 0xbf, 0x08, 0x64, 0x9c, 0x34, 0xa1,
	0x12, 0x06, 0x1d, 0x6b, 0xdd, 0xea, 0xd9, 0x5e, 0x25, 0x0c, 0xc8, 0x06, 0x2c, 0x8f, 0x43, 0x8c,
	0x02, 0xd6, 0xa9, 0xac, 0x5b, 0xbd, 0x5a, 0x7f, 0xc5, 0x55, 0x9c, 0xae, 0xe1, 0x74, 0x77, 0x92,
	0x33, 0x4f, 0x63, 0x9c, 0x13, 0xb8, 0x5b, 0xe0, 0x64, 0x29, 0x4d, 0x18, 0xde, 0x92, 0xf4, 0x31,
	0x90, 0x03, 0x2c, 0x13, 0x9a, 0xb7, 0x2e, 0xa0, 0xfe, 0x49, 0xeb, 0xa7, 0x70, 0x6f, 0x0f, 0x23,
	0xe4, 0x58, 0xd6, 0xbd, 0x07, 0xed, 0xcb, 0xc0, 0x3f, 0x0b, 0x70, 0x7e, 0x54, 0xa0, 0xb6, 0x2b,
	0xa2, 0xcf, 0x86, 0xe9, 0x21, 0xc0, 0xd0, 0xe7, 0xa3, 0xd3, 0x01, 0x0b, 0xcf, 0x51, 0xe2, 0x96,
	0x3c, 0x5b, 0x66, 0x4e, 0xc2, 0x73, 0x24, 0x87, 0xd0, 0x12, 0x69, 0xe0, 0x73, 0x1c, 0x64, 0xaa,
	0x20, 0x17, 0xbe, 0xd0, 0xab, 0xf5, 0xd7, 0x5d, 0x39, 0x7e, 0x77, 0x86, 0xcb, 0xfd, 0x28, 0x91,
	0x3a, 0xf2, 0x9a, 0x62, 0x36, 0x64, 0xdd, 0x9f, 0x16, 0x34, 0x0a, 0x08, 0xd2, 0x86, 0xe5, 0x18,
	0xf9, 0x29, 0x35, 0xfa, 0x74, 0x44, 0xf6, 0xa1, 0x1a, 0xe8, 0x73, 0x68, 0x9b, 0x9e, 0x95, 0x75,
	0x73, 0xa7, 0x07, 0x9f, 0x96, 0x76, 0xdf, 0x40, 0xd5, 0x64, 0x6f, 0x39, 0x87, 0xef, 0x16, 0xd4,
	0x55, 0x6b, 0xed, 0xea, 0x1a, 0xd8, 0xa9, 0xe0, 0x83, 0x11, 0x15, 0x09, 0xd7, 0xa6, 0x55, 0x53,
	0xc1, 0x5f, 0xe7, 0x31, 0x79, 0x02, 0xad, 0x7c, 0x13, 0xb3, 0x8c, 0x66, 0x1a, 0x52, 0x91, 0x90,
	0x46, 0x2a, 0xf8, 0x7e, 0x9e, 0x55, 0xb8, 0x47, 0x50, 0x0f, 0xe4, 0xd0, 0x34, 0x68, 0x41, 0x82,
	0x6a, 0x2a, 0xa7, 0x20, 0x1b, 0x40, 0x94, 0x27, 0x05, 0xb6, 0x45, 0x09, 0xfc, 0x5f, 0xed, 0x5c,
	0x10, 0x3a, 0x6f, 0xa1, 0x71, 0x82, 0x7e, 0x36, 0x3a, 0x35, 0x06, 0xbf, 0x80, 0x26, 0x93, 0x09,
	0x33, 0x3d, 0xa9, 0xf5, 0x6f, 0xa7, 0x6d, 0xb0, 0xd9, 0x62, 0xe7, 0x08, 0x9a, 0x86, 0x4d, 0x9f,
	0x7a, 0x1b, 0x1a, 0x53, 0x3a, 0x26, 0xa2, 0xab, 0xd9, 0xea, 0x86, 0x2d, 0x47, 0x3a, 0xdf, 0x2c,
	0xa8, 0xbe, 0xa7, 0x01, 0x1e, 0x26, 0x63, 0x4a, 0x3a, 0xf0, 0xdf, 0x28, 0x12, 0x8c, 0x63, 0xa6,
	0x27, 0x62, 0x42, 0x42, 0x60, 0x31, 0xa1, 0x01, 0x4a, 0xbf, 0x6c, 0x4f, 0xae, 0xf3, 0x5c, 0x46,
	0x23, 0x94, 0xf6, 0xd8, 0x9e, 0x5c, 0xe7, 0x37, 0x87, 0x71, 0x9f, 0x0b, 0x26, 0xbd, 0xb0, 0x3d,
	0x1d, 0x91, 0x07, 0x60, 0xf3, 0x30, 0x46, 0xc6, 0xfd, 0x38, 0xed, 0x2c, 0xc9, 0xad, 0x8b, 0x84,
	0x13, 0x41, 0xf3, 0x58, 0xf0, 0x5c, 0x86, 0x31, 0xe8, 0x7a, 0x4a, 0x36, 0xc0, 0xce, 0x7f, 0x07,
	0x61, 0x32, 0xa6, 0x52, 0x4e, 0xad, 0xdf, 0xd2, 0x17, 0xd3, 0x9c, 0xcd, 0xab, 0x26, 0x7a, 0xe5,
	0xbc, 0x84, 0xe6, 0x01, 0xde, 0xbc, 0x9b, 0xf3, 0x0a, 0x5a, 0xd3, 0x7a, 0x3d, 0x80, 0x82, 0x00,
	0xab, 0x4c, 0xc0, 0x0e, 0xdc, 0x51, 0x1f, 0x85, 0x9b, 0x6b, 0x88, 0xa1, 0x7d, 0x2c, 0xf8, 0x61,
	0x12, 0xe0, 0xd7, 0x77, 0x7e, 0x9a, 0x86, 0xc9, 0xa4, 0x9c, 0x67, 0x1b, 0x1a, 0x61, 0x5e, 0x30,
	0x88, 0x55, 0xc5, 0x95, 0xff, 0xb0, 0x7a, 0x38, 0xc3, 0xed, 0xf4, 0xa1, 0x7d, 0x80, 0xd7, 0x6b,
	0xe7, 0x7c, 0x80, 0xfb, 0xbf, 0xd5, 0x5c, 0xdc, 0xd7, 0xa2, 0x12, 0x6b, 0x6e, 0x25, 0xcf, 0x61,
	0x55, 0x79, 0x77, 0x3d, 0x31, 0x81, 0x7c, 0x80, 0x54, 0x0d, 0x72, 0xbf, 0xdc, 0xac, 0x2d, 0x00,
	0x2d, 0x11, 0xb9, 0x7f, 0xa5, 0x53, 0x76, 0x68, 0x58, 0x9d, 0x4d, 0xf9, 0xd6, 0xcc, 0xdf, 0xc5,
	0x39, 0x82, 0x95, 0x62, 0x81, 0x36, 0xa8, 0xd8, 0xdd, 0x9a, 0xaf, 0x7b, 0xdf, 0xbc, 0x35, 0xf3,
	0x0b, 0xd8, 0x85, 0x4f, 0x55, 0x43, 0x37, 0x5c, 0x96, 0xab, 0xad, 0x5f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xe2, 0x03, 0x20, 0x64, 0x58, 0x08, 0x00, 0x00,
}